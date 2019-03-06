package healthsvc

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

// Controller ...
type Controller struct {
	cfg    *config.Config
	logger *zap.Logger
}

// NewController returns a new pointer to Controller.
func NewController(c *config.Config) *Controller {
	return &Controller{c, zap.L()}
}

// GetHealth returns the health status of the application.
func (ctl *Controller) GetHealth(c *gin.Context) {
	// reqID := middleware.GetRequestID(c)
	// zap.L().Info(reqID)
	cfg := ctl.cfg

	// Two different way of getting the request id from context.
	ctx := c.Request.Context()

	// Propagate the context to the next layer.
	service(ctx)

	var res Health
	if cfg != nil {
		res = Health{
			BuildDate: cfg.BuildDate,
			GitTag:    cfg.Tag,
			Uptime:    cfg.Uptime(),
		}
	}
	c.JSON(http.StatusOK, res)
}

// GetError simulate an error response.
func (ctl *Controller) GetError(c *gin.Context) {
	ctx := c.Request.Context()
	// Create a logger with the given request id. We can use this to log
	// the request of the endpoint that causes error.
	log := logger.WithContext(ctx)

	// Simulate error.
	err := errors.New("bad error")
	log.Error("endpointError", zap.Error(err))
	model.ErrorJSON(c, err)
}

func service(ctx context.Context) error {
	log := logger.WithContext(ctx)
	log.Info("service: start")
	repository(ctx)
	log.Info("service: end")
	// Stack trace added to this line.
	// return errors.Wrap(errors.New("hello"), "service")
	return nil
}

func repository(ctx context.Context) {
	log := logger.WithContext(ctx)
	log.Info("repository: start")
	// Do work.
	log.Info("repository: end")
}
