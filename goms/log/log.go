package log

import (
	"context"

	"github.com/go-kit/kit/log"
)

type contextKeyType string

const contextLoggerKey contextKeyType = "logger"

func ContextualLogger(ctx context.Context, logger log.Logger, keyvals ...interface{}) (context.Context, log.Logger) {
	if logger := GetLogger(ctx); logger != nil {
		return ctx, logger
	}
	logger = log.With(logger, keyvals...)
	ctx = SetLogger(ctx, logger)
	return ctx, logger
}

func Log(ctx context.Context, keyvals ...interface{}) error {
	return GetLogger(ctx).Log(keyvals...)
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
