package middleware

import (
	"net/http"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func RateLimiter(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		visitor := limiter.GetVisitor(clientIP)
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
