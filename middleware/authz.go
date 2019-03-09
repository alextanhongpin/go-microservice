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
	roleContextKey  = contextKey("role")
	userContextKey  = contextKey("user")
)

func Authz(sign signer.Signer, roles ...string) gin.HandlerFunc {
	roleValidator := set.New()
	for _, role := range roles {
		roleValidator.Add(role)
	}
	checkRole := roleValidator.Size() > 0
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
		role := claims.Role
		if checkRole && !roleValidator.Has(role) {
			err := errors.Errorf(`role "%s" is invalid`, role)
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				model.NewErrorResponse(c, err),
			)
			return
		}
		ctx = ContextWithUser(ctx, user)
		ctx = ContextWithScope(ctx, scope)
		ctx = ContextWithRole(ctx, role)
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

func ContextWithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleContextKey, role)
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleContextKey).(string)
	return role, ok
}
