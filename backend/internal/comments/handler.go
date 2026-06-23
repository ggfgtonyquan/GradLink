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

func (h *Handler) ListComments(c *gin.Context) {
	postID, ok := parseIDParam(c, "id", "Post not found")
	if !ok {
		return
	}

	comments, err := h.Service.ListComments(postID)
	if errors.Is(err, ErrNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(comments))
	for i := range comments {
		response = append(response, commentResponse(&comments[i]))
	}
	c.JSON(http.StatusOK, gin.H{"comments": response})
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

	if IsRejected(comment) {
		auth.AbortWithError(c, http.StatusUnprocessableEntity, "CONTENT_REJECTED", "Content rejected by moderation")
		return
	}
	if IsPendingReview(comment) {
		c.JSON(http.StatusAccepted, commentResponse(comment))
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
	payload := gin.H{
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
	if !comment.CreatedAt.IsZero() {
		payload["createdAt"] = comment.CreatedAt.UTC().Format(http.TimeFormat)
	}
	if comment.Author.Nickname != "" {
		payload["author"] = gin.H{
			"nickname":       comment.Author.Nickname,
			"school":         comment.Author.School,
			"educationLevel": comment.Author.EducationLevel,
			"graduationYear": comment.Author.GraduationYear,
		}
	}
	return payload
}
