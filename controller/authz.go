package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
	"github.com/alextanhongpin/go-microservice/service/authsvc"
)

type Authz struct {
	service authsvc.Service
	signer  signer.Signer
}

func NewAuthz(svc authsvc.Service, sig signer.Signer) *Authz {
	return &Authz{svc, sig}
}

type (
	PostLoginRequest struct {
		authsvc.LoginRequest
	}
	PostLoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (a *Authz) PostLogin(c *gin.Context) {
	var req PostLoginRequest
	if err := c.BindJSON(&req); err != nil {
		model.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := a.service.Login(req.LoginRequest)
	if err != nil {
		log.Error("login user failed", zap.Error(err))
		model.ErrorJSON(c, errors.New("username or password is invalid"))
		return
	}
	var (
		subject = res.Data.ID
		scope   = model.ScopeUser
	)
	token, err := a.service.CreateAccessToken(subject, scope)
	if err != nil {
		log.Error("sign login token failed", zap.Error(err))
		model.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, PostLoginResponse{
		AccessToken: token,
	})
}

type (
	PostRegisterRequest struct {
		authsvc.RegisterRequest
	}
	PostRegisterResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (a *Authz) PostRegister(c *gin.Context) {
	var req PostRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		model.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)

	res, err := a.service.Register(req.RegisterRequest)
	if err != nil {
		log.Error("register user failed", zap.Error(err))
		model.ErrorJSON(c, errors.New("username or email is invalid"))
		return
	}
	var (
		subject = res.Data.ID
		scope   = model.ScopeUser
	)
	token, err := a.service.CreateAccessToken(subject, scope)
	if err != nil {
		log.Error("sign registration token failed", zap.Error(err))
		model.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, PostRegisterResponse{
		AccessToken: token,
	})
}
