package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
	"github.com/alextanhongpin/pkg/set"
)

type contextKey string

func (c contextKey) WithValue(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, c, v)
}

func (c contextKey) Value(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(c).(string)
	return v, ok
}

const (
	// Capitalization matters.
	AuthorizationBearer = "Bearer"
	AuthorizationBasic  = "Basic"

	ScopeContext = contextKey("scope")
	RoleContext  = contextKey("role")
	UserContext  = contextKey("user")
)

func Authz(sign signer.Signer, roles ...api.Role) gin.HandlerFunc {
	roleValidator := set.New()
	for _, role := range roles {
		roleValidator.Add(role)
	}
	checkRole := roleValidator.Size() > 0

	checkAuthorization := func(auth string) (*signer.Claims, error) {
		paths := strings.Split(auth, " ")
		if len(paths) != 2 {
			return nil, errors.New("missing authorization header")
		}
		bearer, token := paths[0], paths[1]
		if valid := strings.EqualFold(bearer, AuthorizationBearer); !valid {
			return nil, errors.New("invalid authorization header")
		}
		claims, err := sign.Verify(token)
		if err != nil {
			return nil, errors.Wrap(err, "middleware verify token failed")
		}

		role := claims.Role
		if checkRole && !roleValidator.Has(api.Role(role)) {
			return nil, errors.Errorf(`role "%s" is invalid`, role)
		}
		return claims, nil
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		claims, err := checkAuthorization(c.GetHeader("Authorization"))
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.NewError(c, err),
			)
			return
		}

		var (
			user  = claims.StandardClaims.Subject
			scope = claims.Scope
			role  = claims.Role
		)

		// Set the context for the next request.
		ctx = UserContext.WithValue(ctx, user)
		ctx = ScopeContext.WithValue(ctx, scope)
		ctx = RoleContext.WithValue(ctx, role)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
