package usersvc

import (
	"net/http"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	UseCase
}

func NewController(usecase UseCase) *Controller {
	return &Controller{usecase}
}

func (ctl *Controller) PostUserInfo(c *gin.Context) {
	type response struct {
		Data User `json:"data"`
	}
	ctx := c.Request.Context()
	id, _ := middleware.UserContext.Value(ctx)
	res, err := ctl.UseCase.UserInfo(id)
	log := logger.WithContext(ctx)
	if err != nil {
		log.Error("get userinfo failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, response{res})
}
