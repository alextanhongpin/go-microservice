package authn

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

type Controller struct {
	UseCase
}

func NewController(usecase UseCase) *Controller {
	return &Controller{usecase}
}

func (ctl *Controller) PostLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := ctl.UseCase.Login(req)
	if err != nil {
		log.Error("login user failed", zap.Error(err))
		api.ErrorJSON(c, errors.New("username or password is invalid"))
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctl *Controller) PostRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := ctl.UseCase.Register(req)
	if err != nil {
		log.Error("register user failed", zap.Error(err))
		api.ErrorJSON(c, errors.New("username or email is invalid"))
		return
	}
	c.JSON(http.StatusOK, res)
}
