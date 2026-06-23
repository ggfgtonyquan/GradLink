package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gradlink/internal/db"
)

type Handler struct {
	Service *Service
}

func NewHandler(database *gorm.DB) *Handler {
	store := &SessionStore{DB: database}
	return &Handler{
		Service: &Service{
			DB:       database,
			Sessions: store,
		},
	}
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid registration payload")
		return
	}

	user, token, err := h.Service.Register(input)
	if errors.Is(err, ErrValidation) {
		AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid registration payload")
		return
	}
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	setSessionCookie(c, token, SessionTTLSeconds())
	c.JSON(http.StatusCreated, userResponse(user))
}

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid login payload")
		return
	}

	user, token, err := h.Service.Login(input.Email, input.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Invalid email or password")
		return
	}
	if errors.Is(err, ErrDisabledUser) {
		AbortWithError(c, http.StatusForbidden, "FORBIDDEN", "User is disabled")
		return
	}
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	setSessionCookie(c, token, SessionTTLSeconds())
	c.JSON(http.StatusOK, userResponse(user))
}

func (h *Handler) Logout(c *gin.Context) {
	cookie, err := c.Cookie(SessionCookieName)
	if err == nil && cookie != "" {
		if err := h.Service.Sessions.Delete(cookie); err != nil {
			AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			return
		}
	}

	setSessionCookie(c, "", -1)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func SessionTTLSeconds() int {
	return int(SessionTTL.Seconds())
}

func setSessionCookie(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SessionCookieName, token, maxAge, "/", "", sessionCookieSecure(), true)
}

func sessionCookieSecure() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("GRADLINK_COOKIE_SECURE")))
	return value == "true" || value == "1" || value == "yes"
}

func userResponse(user *db.User) gin.H {
	return gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
		"role":     user.Role,
		"status":   user.Status,
	}
}
