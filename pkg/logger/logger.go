package logger

import (
	"context"

	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/pkg/reqid"
)

// New returns a basic logger based on the environment.
func New(env, app, hostname string) (logger *zap.Logger) {
	switch env {
	case "production":
		logger, _ = zap.NewProduction()
	case "development":
		logger, _ = zap.NewDevelopment()
	default:
		logger = zap.NewNop()
	}
	logger = logger.With(
		zap.String("app", app),
		zap.String("hostname", hostname),
	)
	return
}

// ReqIDField returns a new logger field for request id.
func ReqIDField(reqID string) zap.Field {
	return zap.String("req_id", reqID)
}

// WithContext creates a new logger and populate the logger with the request
// id.
func WithContext(ctx context.Context) *zap.Logger {
	reqID, _ := reqid.FromContext(ctx)
	if reqID != "" {
		return zap.L().With(ReqIDField(reqID))
	}
	return zap.L()
}
