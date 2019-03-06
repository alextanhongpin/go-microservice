package healthsvc

import (
	"context"
	"net/http"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/xreqid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	cfg *config.Config
}

func NewController(c *config.Config) *Controller {
	return &Controller{c}
}

func (ctl *Controller) GetHealth(c *gin.Context) {
	// reqID := middleware.GetRequestID(c)
	// zap.L().Info(reqID)
	cfg := ctl.cfg

	// Two different way of getting the request id from context.
	ctx := c.Request.Context()
	reqID, _ := xreqid.FromContext(ctx)
	zap.L().Info(reqID)

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

func service(ctx context.Context) error {
	reqID, _ := xreqid.FromContext(ctx)
	log := zap.L().With(logger.ReqIdField(reqID))
	log.Info("service: start")
	repository(ctx)
	log.Info("service: end")
	// Stack trace added to this line.
	// return errors.Wrap(errors.New("hello"), "service")
	return nil
}

func repository(ctx context.Context) {
	reqID, _ := xreqid.FromContext(ctx)
	log := zap.L().With(logger.ReqIdField(reqID))
	log.Info("repository: start")
	// Do work.
	log.Info("repository: end")
}
