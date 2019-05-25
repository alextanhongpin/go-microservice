package authn

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
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
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	var req ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}

	userID, _ := middleware.UserContext.Value(ctx)
	req.ContextUserID = userID

	res, err := ctl.usecase.ChangePassword(ctx, req)
	if err != nil {
		log.Error("update password failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctl *Controller) PostRecoverPassword(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	var req RecoverPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	res, err := ctl.usecase.RecoverPassword(ctx, req)
	if err != nil {
		log.Error("recover password failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	log.Debug("got token", zap.Any("res", res))
	c.JSON(http.StatusOK, res)
}

func (ctl *Controller) PostResetPassword(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	var req ResetPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	res, err := ctl.usecase.ResetPassword(ctx, req)
	if err != nil {
		log.Error("reset password failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctl *Controller) GetResetPasswordView(c *gin.Context) {
	c.HTML(http.StatusOK, "reset_password", nil)
}
