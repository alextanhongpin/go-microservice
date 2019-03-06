package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/pkg/reqid"
)

// Customized version of https://github.com/gin-contrib/zap to include request
// id.
func Logger(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Some evil middlewares modify this values.
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Call the next middleware.
		c.Next()

		ctx := c.Request.Context()
		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}
		reqID, _ := reqid.FromContext(ctx)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(timeFormat)),
			zap.String("request_id", reqID),
			zap.Duration("latency", latency),
		}

		// Include errors if present.
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors[0].Error()))
			// Append error field if this is an erroneous request.
			// for _, e := range c.Errors.Errors() {
			//         logger.Error(e)
			// }
		}

		// Exclude health endpoint, since it introduces noise.
		if path != "/health" {
			logger.Info(path, fields...)
		}
	}
}
