package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, deps *Deps) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
}
