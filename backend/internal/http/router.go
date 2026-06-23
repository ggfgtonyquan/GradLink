package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gradlink/internal/admin"
	"gradlink/internal/auth"
	"gradlink/internal/comments"
	"gradlink/internal/moderation"
	"gradlink/internal/posts"
	"gradlink/internal/reactions"
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
	commentHandler := comments.NewHandler(deps.DB)
	reactionHandler := reactions.NewHandler(deps.DB)
	reportHandler := moderation.NewReportHandler(deps.DB)
	adminHandler := admin.NewHandler(deps.DB)

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/logout", auth.RequireAuth(deps.DB), authHandler.Logout)
	r.GET("/api/tags", tagHandler.List)
	r.GET("/api/posts", postHandler.List)
	r.GET("/api/posts/:id", postHandler.Detail)
	r.GET("/api/posts/:id/comments", commentHandler.ListComments)

	authed := r.Group("/api", auth.RequireAuth(deps.DB))
	authed.GET("/me", profileHandler.Me)
	authed.GET("/me/posts", postHandler.ListMine)
	authed.GET("/me/favorites", postHandler.ListFavorites)
	authed.PUT("/me/profile", profileHandler.PutProfile)
	authed.POST("/posts", auth.RequireProfile(deps.DB), postHandler.Create)
	authed.POST("/posts/:id/comments", auth.RequireProfile(deps.DB), commentHandler.CreateComment)
	authed.POST("/comments/:id/replies", auth.RequireProfile(deps.DB), commentHandler.CreateReply)
	authed.POST("/posts/:id/like", auth.RequireProfile(deps.DB), reactionHandler.Like)
	authed.DELETE("/posts/:id/like", auth.RequireProfile(deps.DB), reactionHandler.Unlike)
	authed.POST("/posts/:id/favorite", auth.RequireProfile(deps.DB), reactionHandler.Favorite)
	authed.DELETE("/posts/:id/favorite", auth.RequireProfile(deps.DB), reactionHandler.Unfavorite)
	authed.POST("/reports", auth.RequireProfile(deps.DB), reportHandler.Create)

	adminRoutes := authed.Group("/admin", admin.RequireAdmin())
	adminRoutes.GET("/moderation", adminHandler.ListModeration)
	adminRoutes.POST("/moderation/:id/approve", adminHandler.ApproveModeration)
	adminRoutes.POST("/moderation/:id/reject", adminHandler.RejectModeration)
	adminRoutes.GET("/reports", adminHandler.ListReports)
	adminRoutes.POST("/reports/:id/resolve", adminHandler.ResolveReport)
	adminRoutes.POST("/reports/:id/reject", adminHandler.RejectReport)
	adminRoutes.GET("/users", adminHandler.ListUsers)
	adminRoutes.POST("/users/:id/ban", adminHandler.BanUser)
	adminRoutes.POST("/users/:id/unban", adminHandler.UnbanUser)
	adminRoutes.GET("/keywords", adminHandler.ListKeywords)
	adminRoutes.POST("/keywords", adminHandler.CreateKeyword)
	adminRoutes.PUT("/keywords/:id", adminHandler.UpdateKeyword)
	adminRoutes.DELETE("/keywords/:id", adminHandler.DeleteKeyword)

	if gin.Mode() == gin.TestMode {
		authed.POST("/test/profile-gated", auth.RequireProfile(deps.DB), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
	}
}
