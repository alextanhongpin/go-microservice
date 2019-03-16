package authn

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

func NewPostLoginController(usecase LoginUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.BindJSON(&req); err != nil {
			api.ErrorJSON(c, err)
			return
		}
		var (
			ctx = c.Request.Context()
			log = logger.WithContext(ctx)
		)
		res, err := usecase(req)
		if err != nil {
			log.Error("login user failed", zap.Error(err))
			api.ErrorJSON(c, errors.New("username or password is invalid"))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func NewPostRegisterController(usecase RegisterUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			api.ErrorJSON(c, err)
			return
		}
		var (
			ctx = c.Request.Context()
			log = logger.WithContext(ctx)
		)
		res, err := usecase(req)
		if err != nil {
			log.Error("register user failed", zap.Error(err))
			api.ErrorJSON(c, errors.New("username or email is invalid"))
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
