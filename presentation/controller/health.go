package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/infrastructure"
)

type Health struct {
	cfg *infrastructure.Config
}

func NewHealth(c *infrastructure.Config) *Health {
	return &Health{c}
}

// The health DTO.
type GetHealthResponse struct {
	BuildDate time.Time `json:"build_date,omitempty"`
	GitTag    string    `json:"git_tag,omitempty"`
	Uptime    string    `json:"uptime"`
}

func (h *Health) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, GetHealthResponse{
		BuildDate: h.cfg.BuildDate,
		GitTag:    h.cfg.Tag,
		Uptime:    h.cfg.Uptime(),
	})
}
