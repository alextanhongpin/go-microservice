package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
	"github.com/alextanhongpin/go-microservice/service/authenticator"
)

type Authenticator struct {
	service authenticator.Service
	signer  signer.Signer
}

func NewAuthenticator(svc authenticator.Service, sig signer.Signer) *Authenticator {
	return &Authenticator{svc, sig}
}

type (
	PostLoginRequest struct {
		authenticator.LoginRequest
	}
	PostLoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (a *Authenticator) PostLogin(c *gin.Context) {
	var req PostLoginRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := a.service.Login(req.LoginRequest)
	if err != nil {
		log.Error("login user failed", zap.Error(err))
		api.ErrorJSON(c, errors.New("username or password is invalid"))
		return
	}
	var (
		subject = res.Data.ID
		scope   = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		role    = api.RoleUser
	)
	token, err := a.service.CreateAccessToken(subject, role.String(), scope)
	if err != nil {
		log.Error("sign login token failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, PostLoginResponse{
		AccessToken: token,
	})
}

type (
	PostRegisterRequest struct {
		authenticator.RegisterRequest
	}
	PostRegisterResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (a *Authenticator) PostRegister(c *gin.Context) {
	var req PostRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)

	res, err := a.service.Register(req.RegisterRequest)
	if err != nil {
		log.Error("register user failed", zap.Error(err))
		api.ErrorJSON(c, errors.New("username or email is invalid"))
		return
	}
	var (
		subject = res.Data.ID
		scope   = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		role    = api.RoleUser
	)
	token, err := a.service.CreateAccessToken(subject, role.String(), scope)
	if err != nil {
		log.Error("sign registration token failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, PostRegisterResponse{
		AccessToken: token,
	})
}
