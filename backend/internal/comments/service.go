package comments

import (
	"errors"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
	"gradlink/internal/db"
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
	Counts          CountView
}

type CountView struct {
	Comments  int
	Likes     int
	Favorites int
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

		created = db.Comment{
			PostID:   post.ID,
			AuthorID: authorID,
			Body:     cleanBody,
			Status:   db.StatusPublished,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
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

		created = db.Comment{
			PostID:          parent.PostID,
			AuthorID:        authorID,
			ParentCommentID: &parent.ID,
			Body:            cleanBody,
			Status:          db.StatusPublished,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
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
		Counts: CountView{
			Comments:  post.CommentCount,
			Likes:     post.LikeCount,
			Favorites: post.FavoriteCount,
		},
	}
}
