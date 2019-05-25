package authn

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

type (
	Controller struct {
		usecase
	}
)

func NewController(u usecase) *Controller {
	return &Controller{u}
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
	accessToken, err := ctl.usecase.LoginWithAccessToken(ctx, req)
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
	accessToken, err := ctl.usecase.RegisterWithAccessToken(ctx, req)
	if err != nil {
		log.Error("register user failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, response{accessToken})
}

func (ctl *Controller) UpdatePassword(c *gin.Context) {
}

func (ctl *Controller) PostRecoverPassword(c *gin.Context) {
}

func (ctl *Controller) PostResetPassword(c *gin.Context) {

}
func (ctl *Controller) GetResetPasswordView(c *gin.Context) {
	c.HTML(http.StatusOK, "reset_password", nil)
}
