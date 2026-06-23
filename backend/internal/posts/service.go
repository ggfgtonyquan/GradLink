package posts

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrNotFound   = errors.New("not found")
)

type Service struct {
	Repository *Repository
}

type CreateInput struct {
	Title  string
	Body   string
	TagIDs []uint
}

type ListInput struct {
	Sort   string
	TagID  *uint
	Query  string
	Limit  int
	Offset int
}

type PostView struct {
	ID        uint
	Title     string
	Body      string
	Excerpt   string
	Status    string
	Tags      []TagView
	Author    AuthorView
	Counts    CountView
	CreatedAt time.Time
}

type TagView struct {
	ID           uint
	Name         string
	Slug         string
	DisplayOrder int
}

type AuthorView struct {
	Nickname       string
	School         any
	EducationLevel string
	GraduationYear int
}

type CountView struct {
	Comments  int
	Likes     int
	Favorites int
}

func (s *Service) CreatePost(authorID uint, input CreateInput) (*PostView, error) {
	title, body, tagIDs, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	tags, err := s.Repository.ListTagsByIDs(tagIDs)
	if err != nil {
		return nil, err
	}
	if len(tags) != len(tagIDs) {
		return nil, ErrValidation
	}

	record, err := s.Repository.CreatePublishedPost(authorID, title, body, tags)
	if err != nil {
		return nil, err
	}

	view := toPostView(record)
	return &view, nil
}

func (s *Service) ListPosts(input ListInput) ([]PostView, error) {
	filter, err := validateListInput(input)
	if err != nil {
		return nil, err
	}

	records, err := s.Repository.ListPublishedPosts(filter)
	if err != nil {
		return nil, err
	}

	views := make([]PostView, 0, len(records))
	for _, record := range records {
		record := record
		views = append(views, toPostView(&record))
	}
	return views, nil
}

func (s *Service) GetPost(postID uint) (*PostView, error) {
	record, err := s.Repository.GetPublishedPost(postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	view := toPostView(record)
	return &view, nil
}

func validateCreateInput(input CreateInput) (string, string, []uint, error) {
	title := strings.TrimSpace(input.Title)
	body := strings.TrimSpace(input.Body)
	if title == "" || body == "" || utf8.RuneCountInString(title) > 160 {
		return "", "", nil, ErrValidation
	}

	if len(input.TagIDs) == 0 || len(input.TagIDs) > 3 {
		return "", "", nil, ErrValidation
	}

	seen := make(map[uint]struct{}, len(input.TagIDs))
	tagIDs := make([]uint, 0, len(input.TagIDs))
	for _, tagID := range input.TagIDs {
		if tagID == 0 {
			return "", "", nil, ErrValidation
		}
		if _, exists := seen[tagID]; exists {
			return "", "", nil, ErrValidation
		}
		seen[tagID] = struct{}{}
		tagIDs = append(tagIDs, tagID)
	}

	return title, body, tagIDs, nil
}

func validateListInput(input ListInput) (FeedFilter, error) {
	sort := strings.TrimSpace(strings.ToLower(input.Sort))
	if sort == "" {
		sort = "latest"
	}
	if sort != "latest" && sort != "hot" {
		return FeedFilter{}, ErrValidation
	}

	limit := input.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 || input.Offset < 0 {
		return FeedFilter{}, ErrValidation
	}

	return FeedFilter{
		Sort:   sort,
		TagID:  input.TagID,
		Query:  strings.TrimSpace(input.Query),
		Limit:  limit,
		Offset: input.Offset,
	}, nil
}

func toPostView(record *PostRecord) PostView {
	tags := make([]TagView, 0, len(record.Tags))
	for _, tag := range record.Tags {
		tags = append(tags, TagView{
			ID:           tag.ID,
			Name:         tag.Name,
			Slug:         tag.Slug,
			DisplayOrder: tag.DisplayOrder,
		})
	}

	return PostView{
		ID:      record.Post.ID,
		Title:   record.Post.Title,
		Body:    record.Post.Body,
		Excerpt: excerptFor(record.Post.Body),
		Status:  record.Post.Status,
		Tags:    tags,
		Author:  authorView(record),
		Counts: CountView{
			Comments:  record.Post.CommentCount,
			Likes:     record.Post.LikeCount,
			Favorites: record.Post.FavoriteCount,
		},
		CreatedAt: record.Post.CreatedAt,
	}
}

func authorView(record *PostRecord) AuthorView {
	view := AuthorView{
		Nickname: record.Post.Author.Nickname,
		School:   nil,
	}
	if record.Profile != nil {
		view.EducationLevel = record.Profile.EducationLevel
		view.GraduationYear = record.Profile.GraduationYear
		if record.Profile.SchoolVisible {
			view.School = record.Profile.School
		}
	}
	return view
}

func excerptFor(body string) string {
	const maxRunes = 120

	runes := []rune(strings.TrimSpace(body))
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes]) + "..."
}
