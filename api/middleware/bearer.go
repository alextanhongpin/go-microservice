package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/passport"
)

const (
	// Capitalization matters.
	AuthorizationBearer = "Bearer"
	AuthorizationBasic  = "Basic"

	ScopeContext = contextKey("scope")
	RoleContext  = contextKey("role")
	UserContext  = contextKey("user")
)

func BearerAuthorizer(sign passport.Signer) gin.HandlerFunc {
	checkAuthorization := func(auth string) (*passport.Claims, error) {
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
