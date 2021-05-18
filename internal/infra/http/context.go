package http

import (
	"context"

	//nolint: depguard
	log "github.com/sirupsen/logrus"
)

type loggerKey struct{}

func ContextWithLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func LoggerFromContext(ctx context.Context) *log.Logger {
	logger, _ := ctx.Value(loggerKey{}).(*log.Logger)
	return logger
}
