package moderation

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gradlink/internal/auth"
	"gradlink/internal/db"
)

var ErrValidation = errors.New("validation failed")

type ReportHandler struct {
	DB *gorm.DB
}

type reportInput struct {
	TargetType string `json:"targetType"`
	TargetID   uint   `json:"targetId"`
	Reason     string `json:"reason"`
	Detail     string `json:"detail"`
}

func NewReportHandler(database *gorm.DB) *ReportHandler {
	return &ReportHandler{DB: database}
}

func (h *ReportHandler) Create(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return
	}

	var input reportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid report payload")
		return
	}

	report, err := h.createReport(user.ID, input)
	if errors.Is(err, ErrValidation) {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid report payload")
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Target not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         report.ID,
		"targetType": report.TargetType,
		"targetId":   report.TargetID,
		"reason":     report.Reason,
		"status":     report.Status,
	})
}

func (h *ReportHandler) createReport(reporterID uint, input reportInput) (*db.Report, error) {
	targetType := strings.TrimSpace(strings.ToLower(input.TargetType))
	reason := strings.TrimSpace(input.Reason)
	detail := strings.TrimSpace(input.Detail)
	if (targetType != "post" && targetType != "comment") || input.TargetID == 0 || reason == "" || utf8.RuneCountInString(reason) > 80 || utf8.RuneCountInString(detail) > 2000 {
		return nil, ErrValidation
	}

	report := db.Report{
		ReporterID: reporterID,
		TargetType: targetType,
		TargetID:   input.TargetID,
		Reason:     reason,
		Detail:     detail,
		Status:     db.StatusPending,
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := validateVisibleTarget(tx, targetType, input.TargetID); err != nil {
			return err
		}
		return tx.Create(&report).Error
	}); err != nil {
		return nil, err
	}
	return &report, nil
}

func validateVisibleTarget(tx *gorm.DB, targetType string, targetID uint) error {
	switch targetType {
	case "post":
		var post db.Post
		return tx.Where("id = ? AND status = ?", targetID, db.StatusPublished).First(&post).Error
	case "comment":
		var comment db.Comment
		return tx.
			Joins("JOIN posts ON posts.id = comments.post_id").
			Where("comments.id = ? AND comments.status = ? AND posts.status = ?", targetID, db.StatusPublished, db.StatusPublished).
			First(&comment).Error
	default:
		return ErrValidation
	}
}
