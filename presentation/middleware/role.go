package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/presentation/api"
)

// RoleChecker checks if the given Actor has the required role.
func RoleChecker(roles ...api.Role) gin.HandlerFunc {
	hasRole := func(role api.Role) bool {
		for _, r := range roles {
			if r == role {
				return true
			}
		}
		return false
	}
	checkRole := len(roles) > 0
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		role, _ := RoleContext.Value(ctx)
		if checkRole && !hasRole(api.Role(role)) {
			err := errors.New("invalid role")
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
