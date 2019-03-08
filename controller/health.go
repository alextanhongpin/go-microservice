package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/config"
)

type Health struct {
	cfg *config.Config
}

// NewHealth returns a new pointer to Controller.
func NewHealth(c *config.Config) *Health {
	return &Health{c}
}

// Health model provides useful information on the app runtime.
type HealthResponse struct {
	BuildDate time.Time `json:"build_date,omitempty"`
	GitTag    string    `json:"git_tag,omitempty"`
	Uptime    string    `json:"uptime"`
}

// GetHealth returns the health status of the application.
func (ctl *Health) GetHealth(c *gin.Context) {
	cfg := ctl.cfg
	c.JSON(http.StatusOK, HealthResponse{
		BuildDate: cfg.BuildDate,
		GitTag:    cfg.Tag,
		Uptime:    cfg.Uptime(),
	})
}
