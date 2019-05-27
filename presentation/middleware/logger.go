package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/pkg/requestid"

	"github.com/alextanhongpin/go-microservice/pkg/logger"
)

// Logger is a customized version of https://github.com/gin-contrib/zap to
// include request id and error.
func Logger(log *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Some evil middlewares modify this values.
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Call the next middleware.
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}
		ctx := c.Request.Context()
		reqID, _ := requestid.Value(ctx)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(timeFormat)),
			zap.Duration("latency", latency),
			logger.ReqIDField(reqID),
		}

		// Include errors if present.
		if len(c.Errors) > 0 {
			fields = append(fields,
				zap.String("error", c.Errors[0].Error()))
		}

		// Exclude health endpoint, since it introduces noise.
		if path != "/health" {
			log.Info(path, fields...)
		}
	}
}
