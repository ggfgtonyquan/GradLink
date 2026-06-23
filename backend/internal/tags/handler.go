package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gradlink/internal/auth"
	"gradlink/internal/db"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(database *gorm.DB) *Handler {
	return &Handler{DB: database}
}

func (h *Handler) List(c *gin.Context) {
	var tagRows []db.Tag
	if err := h.DB.Order("display_order ASC").Order("id ASC").Find(&tagRows).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(tagRows))
	for _, tag := range tagRows {
		response = append(response, gin.H{
			"id":           tag.ID,
			"name":         tag.Name,
			"slug":         tag.Slug,
			"displayOrder": tag.DisplayOrder,
		})
	}

	c.JSON(http.StatusOK, gin.H{"tags": response})
}
