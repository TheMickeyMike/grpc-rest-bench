package main

import (
	"github.com/gin-gonic/gin"
)

func LoadRouter(handler *Handler) *gin.Engine {
	r := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// Logger middleware will write the logs to gin.DefaultWriter when in development mode
	// By default gin.DefaultWriter = os.Stdout
	if gin.Mode() != gin.ReleaseMode {
		r.Use(gin.Logger())
	}

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
		apiv1.GET("/small", handler.SmallJSONResponse)
		apiv1.GET("/users", handler.List)
		apiv1.GET("/users/:id", handler.Get)
		apiv1.DELETE("/users/:id", handler.Delete)
	}

	return r
}
