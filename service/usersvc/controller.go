package usersvc

import (
	"net/http"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/policypalnet/go-pal/log"
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
	id, _ := middleware.UserContext.Value(c.Request.Context())
	res, err := ctl.UseCase.UserInfo(id)
	if err != nil {
		log.Error("get userinfo failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, response{res})
}
