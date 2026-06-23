package reactions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gradlink/internal/auth"
	"gradlink/internal/db"
	"gradlink/internal/posts"
)

const postTargetType = "post"

type Handler struct {
	DB *gorm.DB
}

type counts struct {
	Comments  int
	Likes     int
	Favorites int
}

func NewHandler(database *gorm.DB) *Handler {
	return &Handler{DB: database}
}

func (h *Handler) Like(c *gin.Context) {
	user, postID, ok := currentUserAndPostID(c)
	if !ok {
		return
	}

	post, err := h.setLike(user.ID, postID, true)
	respondCounts(c, post, err)
}

func (h *Handler) Unlike(c *gin.Context) {
	user, postID, ok := currentUserAndPostID(c)
	if !ok {
		return
	}

	post, err := h.setLike(user.ID, postID, false)
	respondCounts(c, post, err)
}

func (h *Handler) Favorite(c *gin.Context) {
	user, postID, ok := currentUserAndPostID(c)
	if !ok {
		return
	}

	post, err := h.setFavorite(user.ID, postID, true)
	respondCounts(c, post, err)
}

func (h *Handler) Unfavorite(c *gin.Context) {
	user, postID, ok := currentUserAndPostID(c)
	if !ok {
		return
	}

	post, err := h.setFavorite(user.ID, postID, false)
	respondCounts(c, post, err)
}

func (h *Handler) setLike(userID, postID uint, liked bool) (*db.Post, error) {
	var post db.Post
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND status = ?", postID, db.StatusPublished).First(&post).Error; err != nil {
			return err
		}

		if liked {
			like := db.Like{UserID: userID, TargetType: postTargetType, TargetID: postID}
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&like)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected > 0 {
				post.LikeCount++
			}
		} else {
			result := tx.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, postTargetType, postID).Delete(&db.Like{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected > 0 && post.LikeCount > 0 {
				post.LikeCount--
			}
		}

		post.HotScore = posts.HotScore(post.CommentCount, post.LikeCount, post.FavoriteCount)
		return tx.Model(&db.Post{}).Where("id = ?", post.ID).Updates(map[string]any{
			"like_count": post.LikeCount,
			"hot_score":  post.HotScore,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (h *Handler) setFavorite(userID, postID uint, favorited bool) (*db.Post, error) {
	var post db.Post
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND status = ?", postID, db.StatusPublished).First(&post).Error; err != nil {
			return err
		}

		if favorited {
			favorite := db.Favorite{UserID: userID, PostID: postID}
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&favorite)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected > 0 {
				post.FavoriteCount++
			}
		} else {
			result := tx.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&db.Favorite{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected > 0 && post.FavoriteCount > 0 {
				post.FavoriteCount--
			}
		}

		post.HotScore = posts.HotScore(post.CommentCount, post.LikeCount, post.FavoriteCount)
		return tx.Model(&db.Post{}).Where("id = ?", post.ID).Updates(map[string]any{
			"favorite_count": post.FavoriteCount,
			"hot_score":      post.HotScore,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func currentUserAndPostID(c *gin.Context) (*db.User, uint, bool) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return nil, 0, false
	}

	parsed, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || parsed == 0 {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return nil, 0, false
	}

	return user, uint(parsed), true
}

func respondCounts(c *gin.Context, post *db.Post, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"postId": post.ID,
		"counts": gin.H{
			"comments":  post.CommentCount,
			"likes":     post.LikeCount,
			"favorites": post.FavoriteCount,
		},
	})
}
