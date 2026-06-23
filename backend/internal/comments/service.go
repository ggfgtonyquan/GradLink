package comments

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
	"gradlink/internal/db"
	"gradlink/internal/moderation"
	"gradlink/internal/posts"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrNotFound   = errors.New("not found")
)

type Service struct {
	DB *gorm.DB
}

type CommentView struct {
	ID              uint
	PostID          uint
	ParentCommentID *uint
	Body            string
	Status          string
	Author          AuthorView
	Counts          CountView
	CreatedAt       time.Time
}

type CountView struct {
	Comments  int
	Likes     int
	Favorites int
}

type AuthorView struct {
	Nickname       string
	School         any
	EducationLevel string
	GraduationYear int
}

func (s *Service) ListComments(postID uint) ([]CommentView, error) {
	var post db.Post
	if err := s.DB.Where("id = ? AND status = ?", postID, db.StatusPublished).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var comments []db.Comment
	if err := s.DB.
		Where("post_id = ? AND status = ?", postID, db.StatusPublished).
		Preload("Author").
		Order("created_at ASC").
		Order("id ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}

	authorIDs := make([]uint, 0, len(comments))
	for _, comment := range comments {
		authorIDs = append(authorIDs, comment.AuthorID)
	}
	profileMap, err := loadProfiles(s.DB, authorIDs)
	if err != nil {
		return nil, err
	}

	views := make([]CommentView, 0, len(comments))
	for _, comment := range comments {
		views = append(views, toCommentViewWithAuthor(comment, post, profileMap[comment.AuthorID]))
	}
	return views, nil
}

func (s *Service) CreateComment(postID, authorID uint, body string) (*CommentView, error) {
	cleanBody, err := validateBody(body)
	if err != nil {
		return nil, err
	}

	var created db.Comment
	var post db.Post
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND status = ?", postID, db.StatusPublished).First(&post).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		decision, err := moderation.ReviewText(tx, cleanBody)
		if err != nil {
			return err
		}
		created = db.Comment{
			PostID:   post.ID,
			AuthorID: authorID,
			Body:     cleanBody,
			Status:   decision.Status,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		if err := moderation.EnqueueIfNeeded(tx, "comment", created.ID, decision); err != nil {
			return err
		}
		if created.Status != db.StatusPublished {
			return tx.First(&post, post.ID).Error
		}

		if err := incrementCommentCount(tx, post.ID); err != nil {
			return err
		}
		return tx.First(&post, post.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return toCommentView(created, post), nil
}

func (s *Service) CreateReply(parentCommentID, authorID uint, body string) (*CommentView, error) {
	cleanBody, err := validateBody(body)
	if err != nil {
		return nil, err
	}

	var created db.Comment
	var post db.Post
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		var parent db.Comment
		err := tx.
			Joins("JOIN posts ON posts.id = comments.post_id").
			Where("comments.id = ? AND comments.status = ? AND posts.status = ?", parentCommentID, db.StatusPublished, db.StatusPublished).
			First(&parent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if parent.ParentCommentID != nil {
			return ErrValidation
		}

		if err := tx.First(&post, parent.PostID).Error; err != nil {
			return err
		}

		decision, err := moderation.ReviewText(tx, cleanBody)
		if err != nil {
			return err
		}
		created = db.Comment{
			PostID:          parent.PostID,
			AuthorID:        authorID,
			ParentCommentID: &parent.ID,
			Body:            cleanBody,
			Status:          decision.Status,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		if err := moderation.EnqueueIfNeeded(tx, "comment", created.ID, decision); err != nil {
			return err
		}
		if created.Status != db.StatusPublished {
			return tx.First(&post, post.ID).Error
		}

		if err := incrementCommentCount(tx, post.ID); err != nil {
			return err
		}
		return tx.First(&post, post.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return toCommentView(created, post), nil
}

func IsPendingReview(comment *CommentView) bool {
	return comment != nil && comment.Status == db.StatusPending
}

func IsRejected(comment *CommentView) bool {
	return comment != nil && comment.Status == db.StatusRejected
}

func validateBody(body string) (string, error) {
	cleanBody := strings.TrimSpace(body)
	if cleanBody == "" || utf8.RuneCountInString(cleanBody) > 2000 {
		return "", ErrValidation
	}
	return cleanBody, nil
}

func incrementCommentCount(tx *gorm.DB, postID uint) error {
	return tx.Model(&db.Post{}).Where("id = ?", postID).Updates(map[string]any{
		"comment_count": gorm.Expr("comment_count + 1"),
		"hot_score":     gorm.Expr("(comment_count + 1) * ? + like_count * ? + favorite_count", posts.CommentHotWeight, posts.LikeHotWeight),
	}).Error
}

func toCommentView(comment db.Comment, post db.Post) *CommentView {
	return &CommentView{
		ID:              comment.ID,
		PostID:          comment.PostID,
		ParentCommentID: comment.ParentCommentID,
		Body:            comment.Body,
		Status:          comment.Status,
		CreatedAt:       comment.CreatedAt,
		Counts: CountView{
			Comments:  post.CommentCount,
			Likes:     post.LikeCount,
			Favorites: post.FavoriteCount,
		},
	}
}

func toCommentViewWithAuthor(comment db.Comment, post db.Post, profile *db.StudentProfile) CommentView {
	view := toCommentView(comment, post)
	view.Author = authorView(comment.Author, profile)
	return *view
}

func authorView(user db.User, profile *db.StudentProfile) AuthorView {
	view := AuthorView{
		Nickname: user.Nickname,
		School:   nil,
	}
	if profile != nil {
		view.EducationLevel = profile.EducationLevel
		view.GraduationYear = profile.GraduationYear
		if profile.SchoolVisible {
			view.School = profile.School
		}
	}
	return view
}

func loadProfiles(database *gorm.DB, authorIDs []uint) (map[uint]*db.StudentProfile, error) {
	if len(authorIDs) == 0 {
		return map[uint]*db.StudentProfile{}, nil
	}

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
