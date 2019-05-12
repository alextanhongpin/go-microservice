package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

type (
	service interface {
		getUsersUseCase
		userInfoUseCase
	}

	Controller struct {
		service
	}
)

func NewController(svc service) *Controller {
	return &Controller{svc}
}

func (ctl *Controller) PostUserInfo(c *gin.Context) {
	type response struct {
		Data *User `json:"data"`
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	id, _ := middleware.UserContext.Value(ctx)
	res, err := ctl.service.UserInfo(id)
	if err != nil {
		log.Error("post userinfo failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, response{res})
}

func (ctl *Controller) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
