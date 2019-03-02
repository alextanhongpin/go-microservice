package logger

import "go.uber.org/zap"

func New(env, name string) (logger *zap.Logger) {
	switch env {
	case "production":
		logger, _ = zap.NewProduction()
	case "development":
		logger, _ = zap.NewDevelopment()
	default:
		logger = zap.NewNop()
	}
	logger = logger.Named(name)
	return
}

func ReqIdField(reqID string) zap.Field {
	return zap.String("req_id", reqID)
}
