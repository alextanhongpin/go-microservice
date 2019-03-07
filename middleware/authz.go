package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
	"github.com/alextanhongpin/pkg/set"
)

type contextKey string

const (
	// Capitalization matters.
	Bearer = "Bearer"
	Basic  = "Basic"

	scopeContextKey = contextKey("scope")
	userContextKey  = contextKey("scope")
)

func Authz(sign signer.Signer, scopes ...string) gin.HandlerFunc {
	s := set.New()
	for _, scope := range scopes {
		s.Add(scope)
	}
	checkScope := s.Size() > 0
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		auth := c.GetHeader("Authorization")
		paths := strings.Split(auth, " ")
		if len(paths) != 2 {
			err := errors.New("missing authorization header")
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				model.NewErrorResponse(c, err),
			)
			return
		}
		bearer, token := paths[0], paths[1]
		if valid := strings.EqualFold(bearer, Bearer); !valid {
			err := errors.New("invalid authorization header")
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				model.NewErrorResponse(c, err),
			)
			return
		}
		claims, err := sign.Verify(token)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				model.NewErrorResponse(c, err),
			)
			return
		}
		user := claims.StandardClaims.Subject
		scope := claims.Scope
		if checkScope && !s.Has(scope) {
			err := errors.Errorf(`scope "%s" is invalid`, scope)
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				model.NewErrorResponse(c, err),
			)
			return
		}
		ctx = ContextWithUser(ctx, user)
		ctx = ContextWithScope(ctx, scope)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func ContextWithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(userContextKey).(string)
	return user, ok
}

func ContextWithScope(ctx context.Context, scope string) context.Context {
	return context.WithValue(ctx, scopeContextKey, scope)
}

func ScopeFromContext(ctx context.Context) (string, bool) {
	scope, ok := ctx.Value(scopeContextKey).(string)
	return scope, ok
}
