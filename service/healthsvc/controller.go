package healthsvc

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/config"
)

// Controller ...
type Controller struct {
	cfg *config.Config
}

// NewController returns a new pointer to Controller.
func NewController(c *config.Config) *Controller {
	return &Controller{c}
}

// Health model provides useful information on the app runtime.
type Health struct {
	BuildDate time.Time `json:"build_date,omitempty"`
	GitTag    string    `json:"git_tag,omitempty"`
	Uptime    string    `json:"uptime"`
}

// GetHealth returns the health status of the application.
func (ctl *Controller) GetHealth(c *gin.Context) {
	cfg := ctl.cfg
	c.JSON(http.StatusOK, Health{
		BuildDate: cfg.BuildDate,
		GitTag:    cfg.Tag,
		Uptime:    cfg.Uptime(),
	})
}
