package posts

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
	return &Handler{
		Service: &Service{
			Repository: &Repository{DB: database},
		},
	}
}

func (h *Handler) Create(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return
	}

	var input struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		TagIDs []uint `json:"tagIds"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid post payload")
		return
	}

	post, err := h.Service.CreatePost(user.ID, CreateInput{
		Title:  input.Title,
		Body:   input.Body,
		TagIDs: input.TagIDs,
	})
	if errors.Is(err, ErrValidation) {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid post payload")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	c.JSON(http.StatusCreated, postResponse(post, true))
}

func (h *Handler) List(c *gin.Context) {
	var tagID *uint
	if raw := c.Query("tagId"); raw != "" {
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || parsed == 0 {
			auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid posts query")
			return
		}
		value := uint(parsed)
		tagID = &value
	}

	views, err := h.Service.ListPosts(ListInput{
		Sort:  c.Query("sort"),
		TagID: tagID,
		Query: c.Query("q"),
	})
	if errors.Is(err, ErrValidation) {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid posts query")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(views))
	for i := range views {
		response = append(response, postResponse(&views[i], false))
	}
	c.JSON(http.StatusOK, gin.H{"posts": response})
}

func (h *Handler) Detail(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || postID == 0 {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}

	post, err := h.Service.GetPost(uint(postID))
	if errors.Is(err, ErrNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	c.JSON(http.StatusOK, postResponse(post, true))
}

func postResponse(post *PostView, includeBody bool) gin.H {
	tags := make([]gin.H, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, gin.H{
			"id":           tag.ID,
			"name":         tag.Name,
			"slug":         tag.Slug,
			"displayOrder": tag.DisplayOrder,
		})
	}

	payload := gin.H{
		"id":        post.ID,
		"title":     post.Title,
		"excerpt":   post.Excerpt,
		"status":    post.Status,
		"tags":      tags,
		"author":    authorResponse(post.Author),
		"counts":    countsResponse(post.Counts),
		"createdAt": post.CreatedAt.UTC().Format(http.TimeFormat),
	}
	if includeBody {
		payload["body"] = post.Body
	}

	return payload
}

func authorResponse(author AuthorView) gin.H {
	return gin.H{
		"nickname":       author.Nickname,
		"school":         author.School,
		"educationLevel": author.EducationLevel,
		"graduationYear": author.GraduationYear,
	}
}

func countsResponse(counts CountView) gin.H {
	return gin.H{
		"comments":  counts.Comments,
		"likes":     counts.Likes,
		"favorites": counts.Favorites,
	}
}
