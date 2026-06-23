package comments

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gradlink/internal/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(database *gorm.DB) *Handler {
	return &Handler{Service: &Service{DB: database}}
}

func (h *Handler) CreateComment(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return
	}

	postID, ok := parseIDParam(c, "id", "Post not found")
	if !ok {
		return
	}

	var input struct {
		Body string `json:"body"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid comment payload")
		return
	}

	comment, err := h.Service.CreateComment(postID, user.ID, input.Body)
	h.respondComment(c, comment, err)
}

func (h *Handler) CreateReply(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return
	}

	parentCommentID, ok := parseIDParam(c, "id", "Comment not found")
	if !ok {
		return
	}

	var input struct {
		Body string `json:"body"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid reply payload")
		return
	}

	comment, err := h.Service.CreateReply(parentCommentID, user.ID, input.Body)
	h.respondComment(c, comment, err)
}

func (h *Handler) respondComment(c *gin.Context, comment *CommentView, err error) {
	if errors.Is(err, ErrValidation) {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid comment payload")
		return
	}
	if errors.Is(err, ErrNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Target not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	c.JSON(http.StatusCreated, commentResponse(comment))
}

func parseIDParam(c *gin.Context, name, message string) (uint, bool) {
	parsed, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || parsed == 0 {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", message)
		return 0, false
	}
	return uint(parsed), true
}

func commentResponse(comment *CommentView) gin.H {
	return gin.H{
		"id":              comment.ID,
		"postId":          comment.PostID,
		"parentCommentId": comment.ParentCommentID,
		"body":            comment.Body,
		"status":          comment.Status,
		"counts": gin.H{
			"comments":  comment.Counts.Comments,
			"likes":     comment.Counts.Likes,
			"favorites": comment.Counts.Favorites,
		},
	}
}
