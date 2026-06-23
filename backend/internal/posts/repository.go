package posts

import (
	"strings"

	"gorm.io/gorm"
	"gradlink/internal/db"
	"gradlink/internal/moderation"
)

type Repository struct {
	DB *gorm.DB
}

type FeedFilter struct {
	Sort   string
	TagID  *uint
	Query  string
	Limit  int
	Offset int
}

type PostRecord struct {
	Post    db.Post
	Tags    []db.Tag
	Profile *db.StudentProfile
}

func (r *Repository) ListTagsByIDs(tagIDs []uint) ([]db.Tag, error) {
	var tags []db.Tag
	err := r.DB.
		Where("id IN ?", tagIDs).
		Order("display_order ASC").
		Order("id ASC").
		Find(&tags).Error
	return tags, err
}

func (r *Repository) CreatePost(authorID uint, title, body string, tags []db.Tag) (*PostRecord, error) {
	var postID uint
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		decision, err := moderation.ReviewText(tx, title+"\n"+body)
		if err != nil {
			return err
		}

		post := db.Post{
			AuthorID: authorID,
			Title:    title,
			Body:     body,
			Status:   decision.Status,
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}

		postTags := make([]db.PostTag, 0, len(tags))
		for _, tag := range tags {
			postTags = append(postTags, db.PostTag{
				PostID: post.ID,
				TagID:  tag.ID,
			})
		}
		if err := tx.Create(&postTags).Error; err != nil {
			return err
		}
		if err := moderation.EnqueueIfNeeded(tx, "post", post.ID, decision); err != nil {
			return err
		}

		postID = post.ID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetPost(postID)
}

func (r *Repository) ListPublishedPosts(filter FeedFilter) ([]PostRecord, error) {
	query := r.DB.Model(&db.Post{}).
		Where("posts.status = ?", db.StatusPublished)

	if filter.TagID != nil {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id = ?", *filter.TagID)
	}
	if filter.Query != "" {
		like := "%" + escapeLikePattern(filter.Query) + "%"
		query = query.Where("(posts.title LIKE ? ESCAPE '\\' OR posts.body LIKE ? ESCAPE '\\')", like, like)
	}

	switch filter.Sort {
	case "hot":
		query = query.Order("posts.hot_score DESC").Order("posts.created_at DESC").Order("posts.id DESC")
	default:
		query = query.Order("posts.created_at DESC").Order("posts.id DESC")
	}

	var postRows []db.Post
	if err := query.
		Distinct("posts.id", "posts.author_id", "posts.title", "posts.body", "posts.status", "posts.comment_count", "posts.like_count", "posts.favorite_count", "posts.hot_score", "posts.created_at", "posts.updated_at").
		Preload("Author").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&postRows).Error; err != nil {
		return nil, err
	}

	return r.loadPostRecords(postRows)
}

func (r *Repository) ListPostsByAuthor(authorID uint, filter FeedFilter) ([]PostRecord, error) {
	var postRows []db.Post
	if err := r.DB.Model(&db.Post{}).
		Where("posts.author_id = ?", authorID).
		Preload("Author").
		Order("posts.created_at DESC").
		Order("posts.id DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&postRows).Error; err != nil {
		return nil, err
	}

	return r.loadPostRecords(postRows)
}

func (r *Repository) ListFavoritePosts(userID uint, filter FeedFilter) ([]PostRecord, error) {
	var postRows []db.Post
	if err := r.DB.Model(&db.Post{}).
		Joins("JOIN favorites ON favorites.post_id = posts.id").
		Where("favorites.user_id = ? AND posts.status = ?", userID, db.StatusPublished).
		Preload("Author").
		Order("favorites.created_at DESC").
		Order("posts.id DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&postRows).Error; err != nil {
		return nil, err
	}

	return r.loadPostRecords(postRows)
}

func (r *Repository) GetPublishedPost(postID uint) (*PostRecord, error) {
	var post db.Post
	err := r.DB.
		Where("posts.id = ? AND posts.status = ?", postID, db.StatusPublished).
		Preload("Author").
		First(&post).Error
	if err != nil {
		return nil, err
	}

	records, err := r.loadPostRecords([]db.Post{post})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &records[0], nil
}

func (r *Repository) GetPost(postID uint) (*PostRecord, error) {
	var post db.Post
	err := r.DB.
		Where("posts.id = ?", postID).
		Preload("Author").
		First(&post).Error
	if err != nil {
		return nil, err
	}

	records, err := r.loadPostRecords([]db.Post{post})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &records[0], nil
}

func (r *Repository) loadPostRecords(postRows []db.Post) ([]PostRecord, error) {
	if len(postRows) == 0 {
		return []PostRecord{}, nil
	}

	postIDs := make([]uint, 0, len(postRows))
	authorIDs := make([]uint, 0, len(postRows))
	for _, post := range postRows {
		postIDs = append(postIDs, post.ID)
		authorIDs = append(authorIDs, post.AuthorID)
	}

	profileMap, err := loadProfiles(r.DB, authorIDs)
	if err != nil {
		return nil, err
	}
	tagMap, err := loadTagsByPostID(r.DB, postIDs)
	if err != nil {
		return nil, err
	}

	records := make([]PostRecord, 0, len(postRows))
	for _, post := range postRows {
		records = append(records, PostRecord{
			Post:    post,
			Tags:    tagMap[post.ID],
			Profile: profileMap[post.AuthorID],
		})
	}
	return records, nil
}

func escapeLikePattern(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}

func loadProfiles(database *gorm.DB, authorIDs []uint) (map[uint]*db.StudentProfile, error) {
	var profiles []db.StudentProfile
	if err := database.Where("user_id IN ?", authorIDs).Find(&profiles).Error; err != nil {
		return nil, err
	}

	profileMap := make(map[uint]*db.StudentProfile, len(profiles))
	for i := range profiles {
		profile := profiles[i]
		profileMap[profile.UserID] = &profile
	}
	return profileMap, nil
}

func loadTagsByPostID(database *gorm.DB, postIDs []uint) (map[uint][]db.Tag, error) {
	type row struct {
		PostID       uint
		ID           uint
		Name         string
		Slug         string
		DisplayOrder int
	}

	var rows []row
	if err := database.
		Table("tags").
		Select("post_tags.post_id, tags.id, tags.name, tags.slug, tags.display_order").
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id IN ?", postIDs).
		Order("post_tags.post_id ASC").
		Order("tags.display_order ASC").
		Order("tags.id ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]db.Tag, len(postIDs))
	for _, row := range rows {
		tagMap[row.PostID] = append(tagMap[row.PostID], db.Tag{
			ID:           row.ID,
			Name:         row.Name,
			Slug:         row.Slug,
			DisplayOrder: row.DisplayOrder,
		})
	}
	return tagMap, nil
}
