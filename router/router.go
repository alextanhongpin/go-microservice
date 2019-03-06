package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/middleware"
	"github.com/alextanhongpin/go-microservice/service/healthsvc"
)

// New returns a new Router.
func New(cfg *config.Config) http.Handler {
	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Custom middlewares.
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(zap.L(), time.RFC3339, true))
	// TODO: Include authorization signer.

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
