package authn

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/passport"
)

type Controller struct {
	service  Service
	passport passport.Signer
}

func NewController(svc Service, sig passport.Signer) *Controller {
	return &Controller{svc, sig}
}

type (
	PostLoginRequest struct {
		LoginRequest
	}
	PostLoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (ctl *Controller) PostLogin(c *gin.Context) {
	var req PostLoginRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)
	res, err := ctl.service.Login(req.LoginRequest)
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
	token, err := ctl.service.CreateAccessToken(subject, role.String(), scope)
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
		RegisterRequest
	}
	PostRegisterResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (ctl *Controller) PostRegister(c *gin.Context) {
	var req PostRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		api.ErrorJSON(c, err)
		return
	}
	var (
		ctx = c.Request.Context()
		log = logger.WithContext(ctx)
	)

	res, err := ctl.service.Register(req.RegisterRequest)
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
	token, err := ctl.service.CreateAccessToken(subject, role.String(), scope)
	if err != nil {
		log.Error("sign registration token failed", zap.Error(err))
		api.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, PostRegisterResponse{
		AccessToken: token,
	})
}
