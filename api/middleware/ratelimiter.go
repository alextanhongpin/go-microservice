package middleware

import (
	"fmt"
	"net/http"

	"github.com/alextanhongpin/go-microservice/api"
	ratelimit "github.com/alextanhongpin/pkg/ratelimiter"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
