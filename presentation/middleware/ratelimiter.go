package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/presentation/api"
	ratelimit "github.com/alextanhongpin/pkg/ratelimiter"
)

func RateLimiter(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		// For per path, consider concatenating the c.Request.URL.Path
		// with the client IP.
		visitor := limiter.GetVisitor(fmt.Sprintf("%s/%s", c.Request.URL.Path, clientIP))
		if !visitor.Allow() {
			err := errors.Errorf(`client ip "%s" has too many requests`, clientIP)
			c.Error(err)
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				api.NewError(c, err),
			)
			return
		}
		c.Next()
	}
}
