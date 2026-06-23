package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gradlink/internal/db"
)

const currentUserKey = "currentUser"

func RequireAuth(database *gorm.DB) gin.HandlerFunc {
	store := &SessionStore{DB: database}

	return func(c *gin.Context) {
		cookie, err := c.Cookie(SessionCookieName)
		if err != nil || cookie == "" {
			AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
			return
		}

		session, err := store.FindActive(cookie)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
			return
		}
		if err != nil {
			AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			return
		}
		if session.User.Status != db.UserStatusActive {
			AbortWithError(c, http.StatusForbidden, "FORBIDDEN", "User is disabled")
			return
		}

		c.Set(currentUserKey, &session.User)
		c.Next()
	}
}

func RequireProfile(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok {
			AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
			return
		}

		var count int64
		if err := database.Model(&db.StudentProfile{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
			AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			return
		}
		if count == 0 {
			AbortWithError(c, http.StatusForbidden, "PROFILE_REQUIRED", "Complete your profile before continuing")
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*db.User, bool) {
	value, ok := c.Get(currentUserKey)
	if !ok {
		return nil, false
	}

	user, ok := value.(*db.User)
	return user, ok
}

func AbortWithError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
