package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/presentation/api"
	"github.com/alextanhongpin/go-microservice/presentation/middleware"
	"github.com/alextanhongpin/go-microservice/usecase"
)

type (
	Authn struct {
		usecase usecase.Authn
	}
)

func NewAuthn(u usecase.Authn) *Authn {
	return &Authn{u}
}

func (a *Authn) PostLogin(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := a.usecase.Login(ctx, req)
	if err != nil {
		log.Error("login user failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (a *Authn) PostRegister(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := a.usecase.Register(ctx, req)
	if err != nil {
		log.Error("register user failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (a *Authn) UpdatePassword(c *gin.Context) {
	var req usecase.ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)

	userID, _ := middleware.UserContext.Value(ctx)
	req.ContextUserID = userID

	res, err := a.usecase.ChangePassword(ctx, req)
	if err != nil {
		log.Error("update password failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (a *Authn) PostRecoverPassword(c *gin.Context) {
	var req usecase.RecoverPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := a.usecase.RecoverPassword(ctx, req)
	if err != nil {
		log.Error("recover password failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (a *Authn) PostResetPassword(c *gin.Context) {
	var req usecase.ResetPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := a.usecase.ResetPassword(ctx, req)
	if err != nil {
		log.Error("reset password failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// func (ctl *Authn) GetResetPasswordView(c *gin.Context) {
//         c.HTML(http.StatusOK, "reset_password", nil)
// }
