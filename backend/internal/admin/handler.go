package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.CurrentUser(c)
		if !ok {
			auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
			return
		}
		if user.Role != db.UserRoleAdmin {
			auth.AbortWithError(c, http.StatusForbidden, "FORBIDDEN", "Admin access required")
			return
		}

		c.Next()
	}
}

func (h *Handler) ListModeration(c *gin.Context) {
	var items []db.ModerationQueueItem
	if err := h.DB.Where("status = ?", db.StatusPending).Order("created_at ASC").Order("id ASC").Find(&items).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(items))
	for _, item := range items {
		payload, err := h.moderationItemWithContent(item)
		if err != nil {
			auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			return
		}
		response = append(response, payload)
	}
	c.JSON(http.StatusOK, gin.H{"items": response})
}

func (h *Handler) ApproveModeration(c *gin.Context) {
	h.decideModeration(c, db.StatusPublished)
}

func (h *Handler) RejectModeration(c *gin.Context) {
	h.decideModeration(c, db.StatusRejected)
}

func (h *Handler) ListReports(c *gin.Context) {
	var reports []db.Report
	if err := h.DB.Order("created_at ASC").Order("id ASC").Find(&reports).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(reports))
	for _, report := range reports {
		response = append(response, reportResponse(report))
	}
	c.JSON(http.StatusOK, gin.H{"reports": response})
}

func (h *Handler) ResolveReport(c *gin.Context) {
	h.updateReportStatus(c, db.ReportStatusResolved)
}

func (h *Handler) RejectReport(c *gin.Context) {
	h.updateReportStatus(c, db.StatusRejected)
}

func (h *Handler) updateReportStatus(c *gin.Context, status string) {
	reportID, ok := parseIDParam(c, "id")
	if !ok {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Report not found")
		return
	}

	var report db.Report
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&report, reportID).Error; err != nil {
			return err
		}
		return tx.Model(&report).Update("status", status).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Report not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	if err := h.DB.First(&report, reportID).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	c.JSON(http.StatusOK, reportResponse(report))
}

func (h *Handler) ListUsers(c *gin.Context) {
	var users []db.User
	if err := h.DB.Order("created_at ASC").Order("id ASC").Find(&users).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(users))
	for _, user := range users {
		response = append(response, userResponse(user))
	}
	c.JSON(http.StatusOK, gin.H{"users": response})
}

func (h *Handler) BanUser(c *gin.Context) {
	h.updateUserStatus(c, db.UserStatusDisabled)
}

func (h *Handler) UnbanUser(c *gin.Context) {
	h.updateUserStatus(c, db.UserStatusActive)
}

func (h *Handler) updateUserStatus(c *gin.Context, status string) {
	userID, ok := parseIDParam(c, "id")
	if !ok {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	var user db.User
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		return tx.Model(&user).Update("status", status).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	if err := h.DB.First(&user, userID).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	c.JSON(http.StatusOK, userResponse(user))
}

func (h *Handler) ListKeywords(c *gin.Context) {
	var keywords []db.ModerationKeyword
	if err := h.DB.Order("created_at ASC").Order("id ASC").Find(&keywords).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := make([]gin.H, 0, len(keywords))
	for _, keyword := range keywords {
		response = append(response, keywordResponse(keyword))
	}
	c.JSON(http.StatusOK, gin.H{"keywords": response})
}

func (h *Handler) CreateKeyword(c *gin.Context) {
	input, ok := bindKeywordInput(c)
	if !ok {
		return
	}

	keyword := db.ModerationKeyword{
		Keyword:   input.Keyword,
		RiskLevel: input.RiskLevel,
		Action:    input.Action,
		Enabled:   *input.Enabled,
	}
	if err := h.DB.Create(&keyword).Error; err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid keyword payload")
		return
	}

	c.JSON(http.StatusCreated, keywordResponse(keyword))
}

func (h *Handler) UpdateKeyword(c *gin.Context) {
	keywordID, ok := parseIDParam(c, "id")
	if !ok {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Keyword not found")
		return
	}
	input, ok := bindKeywordInput(c)
	if !ok {
		return
	}

	var keyword db.ModerationKeyword
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&keyword, keywordID).Error; err != nil {
			return err
		}
		return tx.Model(&keyword).Updates(map[string]any{
			"keyword":    input.Keyword,
			"risk_level": input.RiskLevel,
			"action":     input.Action,
			"enabled":    *input.Enabled,
		}).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Keyword not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid keyword payload")
		return
	}

	if err := h.DB.First(&keyword, keywordID).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	c.JSON(http.StatusOK, keywordResponse(keyword))
}

func (h *Handler) DeleteKeyword(c *gin.Context) {
	keywordID, ok := parseIDParam(c, "id")
	if !ok {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Keyword not found")
		return
	}

	result := h.DB.Delete(&db.ModerationKeyword{}, keywordID)
	if result.Error != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	if result.RowsAffected == 0 {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Keyword not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) decideModeration(c *gin.Context, status string) {
	user, _ := auth.CurrentUser(c)
	itemID, ok := parseIDParam(c, "id")
	if !ok {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Moderation item not found")
		return
	}

	var item db.ModerationQueueItem
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&item, itemID).Error; err != nil {
			return err
		}
		if err := updateTargetStatus(tx, item.TargetType, item.TargetID, status); err != nil {
			return err
		}
		return tx.Model(&item).Updates(map[string]any{
			"status":      status,
			"reviewer_id": user.ID,
			"decision":    status,
		}).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		auth.AbortWithError(c, http.StatusNotFound, "NOT_FOUND", "Moderation item not found")
		return
	}
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	if err := h.DB.First(&item, itemID).Error; err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	c.JSON(http.StatusOK, moderationItemResponse(item))
}

func updateTargetStatus(tx *gorm.DB, targetType string, targetID uint, status string) error {
	switch targetType {
	case "post":
		result := tx.Model(&db.Post{}).Where("id = ?", targetID).Update("status", status)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	case "comment":
		var comment db.Comment
		if err := tx.First(&comment, targetID).Error; err != nil {
			return err
		}
		previousStatus := comment.Status
		if err := tx.Model(&comment).Update("status", status).Error; err != nil {
			return err
		}
		if previousStatus != db.StatusPublished && status == db.StatusPublished {
			return incrementPostCommentCount(tx, comment.PostID)
		}
		return nil
	default:
		return errors.New("unsupported moderation target")
	}
}

func incrementPostCommentCount(tx *gorm.DB, postID uint) error {
	return tx.Model(&db.Post{}).Where("id = ?", postID).Updates(map[string]any{
		"comment_count": gorm.Expr("comment_count + 1"),
		"hot_score":     gorm.Expr("(comment_count + 1) * 3 + like_count * 2 + favorite_count"),
	}).Error
}

func parseIDParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		return 0, false
	}
	return uint(value), true
}

func moderationItemResponse(item db.ModerationQueueItem) gin.H {
	return gin.H{
		"id":              item.ID,
		"targetType":      item.TargetType,
		"targetId":        item.TargetID,
		"matchedKeywords": item.MatchedKeywords,
		"status":          item.Status,
		"reviewerId":      item.ReviewerID,
		"decision":        item.Decision,
	}
}

func (h *Handler) moderationItemWithContent(item db.ModerationQueueItem) (gin.H, error) {
	payload := moderationItemResponse(item)
	switch item.TargetType {
	case "post":
		var post db.Post
		if err := h.DB.Preload("Author").First(&post, item.TargetID).Error; err != nil {
			return nil, err
		}
		payload["content"] = gin.H{
			"title":    post.Title,
			"body":     post.Body,
			"authorId": post.AuthorID,
			"author":   post.Author.Nickname,
			"status":   post.Status,
		}
	case "comment":
		var comment db.Comment
		if err := h.DB.Preload("Author").First(&comment, item.TargetID).Error; err != nil {
			return nil, err
		}
		payload["content"] = gin.H{
			"title":    "评论",
			"body":     comment.Body,
			"postId":   comment.PostID,
			"authorId": comment.AuthorID,
			"author":   comment.Author.Nickname,
			"status":   comment.Status,
		}
	}
	return payload, nil
}

func reportResponse(report db.Report) gin.H {
	return gin.H{
		"id":         report.ID,
		"targetType": report.TargetType,
		"targetId":   report.TargetID,
		"reason":     report.Reason,
		"detail":     report.Detail,
		"status":     report.Status,
		"reporterId": report.ReporterID,
	}
}

func userResponse(user db.User) gin.H {
	return gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
		"role":     user.Role,
		"status":   user.Status,
	}
}

type keywordInput struct {
	Keyword   string `json:"keyword"`
	RiskLevel string `json:"riskLevel"`
	Action    string `json:"action"`
	Enabled   *bool  `json:"enabled"`
}

func bindKeywordInput(c *gin.Context) (keywordInput, bool) {
	var raw keywordInput
	if err := c.ShouldBindJSON(&raw); err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid keyword payload")
		return keywordInput{}, false
	}

	raw.Keyword = strings.TrimSpace(raw.Keyword)
	raw.RiskLevel = strings.TrimSpace(raw.RiskLevel)
	raw.Action = strings.TrimSpace(raw.Action)
	if raw.Enabled == nil {
		defaultEnabled := true
		raw.Enabled = &defaultEnabled
	}
	if raw.Keyword == "" || raw.RiskLevel == "" || (raw.Action != db.StatusPending && raw.Action != db.StatusRejected) {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid keyword payload")
		return keywordInput{}, false
	}

	return raw, true
}

func keywordResponse(keyword db.ModerationKeyword) gin.H {
	return gin.H{
		"id":        keyword.ID,
		"keyword":   keyword.Keyword,
		"riskLevel": keyword.RiskLevel,
		"action":    keyword.Action,
		"enabled":   keyword.Enabled,
	}
}
