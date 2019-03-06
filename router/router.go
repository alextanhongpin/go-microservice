package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/middleware"
	"github.com/alextanhongpin/go-microservice/service/healthsvc"
)

func New(cfg *config.Config) http.Handler {
	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	// TODO: Include cors.

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
