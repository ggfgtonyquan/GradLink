package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gradlink/internal/auth"
	"gradlink/internal/posts"
	"gradlink/internal/tags"
	"gradlink/internal/users"
)

func RegisterRoutes(r *gin.Engine, deps *Deps) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	if deps == nil || deps.DB == nil {
		return
	}

	authHandler := auth.NewHandler(deps.DB)
	profileHandler := users.NewProfileHandler(deps.DB)
	postHandler := posts.NewHandler(deps.DB)
	tagHandler := tags.NewHandler(deps.DB)

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/logout", auth.RequireAuth(deps.DB), authHandler.Logout)
	r.GET("/api/tags", tagHandler.List)
	r.GET("/api/posts", postHandler.List)
	r.GET("/api/posts/:id", postHandler.Detail)

	authed := r.Group("/api", auth.RequireAuth(deps.DB))
	authed.GET("/me", profileHandler.Me)
	authed.PUT("/me/profile", profileHandler.PutProfile)
	authed.POST("/posts", auth.RequireProfile(deps.DB), postHandler.Create)

	if gin.Mode() == gin.TestMode {
		authed.POST("/test/profile-gated", auth.RequireProfile(deps.DB), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
	}
}
