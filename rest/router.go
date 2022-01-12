package main

import (
	"github.com/gin-gonic/gin"
)

func LoadRouter(handler *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	status := r.Group("/status")
	{
		status.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/users", handler.List)
		apiv1.GET("/users/:id", handler.Get)
		apiv1.DELETE("/users/:id", handler.Delete)
	}

	return r
}
