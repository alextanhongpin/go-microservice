package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/middleware"
	"github.com/alextanhongpin/go-microservice/service/healthsvc"
)

func New(cfg *config.Config) http.Handler {
	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(zap.L(), time.RFC3339, true))
	// TODO: Include cors.
	// TODO: Include logger, but exclude the /health path.

	// Health endpoint.
	{
		ctl := healthsvc.NewController(cfg)
		r.GET("/health", ctl.GetHealth)
		r.GET("/error", ctl.GetError)
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
