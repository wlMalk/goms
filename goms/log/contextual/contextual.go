package contextual

import (
	"context"

	"github.com/go-kit/kit/log"
)

type contextKeyType string

const contextLoggerKey contextKeyType = "logger"

func ContextualLogger(ctx context.Context, logger log.Logger, keyvals ...interface{}) log.Logger {
	if logger := GetLogger(ctx); logger != nil {
		return logger
	}
	logger = log.With(logger, keyvals...)
	return logger
}

func SetLogger(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, logger)
}

func GetLogger(ctx context.Context) log.Logger {
	logger := ctx.Value(contextLoggerKey)
	if logger == nil {
		return nil
	}
	return logger.(log.Logger)
}
