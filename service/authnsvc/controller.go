package authnsvc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

type (
	Controller struct {
		service
	}
)

func NewController(svc service) *Controller {
	return &Controller{svc}
}

func (ctl *Controller) PostLogin(c *gin.Context) {
	type request = LoginRequest
	type response struct {
		AccessToken string `json:"access_token"`
	}
	var req request
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	accessToken, err := ctl.service.LoginWithAccessToken(ctx, req)
	if err != nil {
		log.Error("login user failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, response{accessToken})
}

func (ctl *Controller) PostRegister(c *gin.Context) {
	type request = RegisterRequest
	type response struct {
		AccessToken string `json:"access_token"`
	}
	var req request
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	accessToken, err := ctl.service.RegisterWithAccessToken(ctx, req)
	if err != nil {
		log.Error("register user failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, response{accessToken})
}
