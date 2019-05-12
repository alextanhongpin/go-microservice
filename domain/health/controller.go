package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/infrastructure"
)

type Controller struct {
	cfg *infrastructure.Config
}

// NewController returns a new pointer to Controller.
func NewController(c *infrastructure.Config) *Controller {
	return &Controller{c}
}

// GetController returns the health status of the application.
func (ctl *Controller) GetHealth(c *gin.Context) {
	type response struct {
		BuildDate time.Time `json:"build_date,omitempty"`
		GitTag    string    `json:"git_tag,omitempty"`
		Uptime    string    `json:"uptime"`
	}
	c.JSON(http.StatusOK, response{
		BuildDate: ctl.cfg.BuildDate,
		GitTag:    ctl.cfg.Tag,
		Uptime:    ctl.cfg.Uptime(),
	})
}
