package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
)

func Scope(scope api.Scope) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		s, _ := ScopeContext.Value(ctx)
		if !scope.Equal(s) {
			err := errors.New("invalid scope")
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.NewError(c, err),
			)
			return
		}
		c.Next()

	}
}
