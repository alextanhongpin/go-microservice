package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/logging/config"
	"github.com/alextanhongpin/logging/middleware"
	"github.com/alextanhongpin/logging/service/healthsvc"
)

func New(cfg *config.Config) http.Handler {
	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())

	// Health endpoint.
	{
		ctl := healthsvc.NewController(cfg)
		r.GET("/health", ctl.GetHealth)
	}

	// Handle no route.
	r.NoRoute(func(c *gin.Context) {
		// TODO: Cleanup message.
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})
	return r
}
