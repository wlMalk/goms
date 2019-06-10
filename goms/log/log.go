package log

import (
	"context"

	"github.com/wlMalk/goms/goms/log/contextual"

	"github.com/go-kit/kit/log/level"
)

func Log(ctx context.Context, keyvals ...interface{}) error {
	return contextual.GetLogger(ctx).Log(keyvals...)
}

func Error(ctx context.Context, keyvals ...interface{}) error {
	return level.Error(contextual.GetLogger(ctx)).Log(keyvals...)
}

func Warn(ctx context.Context, keyvals ...interface{}) error {
	return level.Warn(contextual.GetLogger(ctx)).Log(keyvals...)
}

func Info(ctx context.Context, keyvals ...interface{}) error {
	return level.Info(contextual.GetLogger(ctx)).Log(keyvals...)
}

func Debug(ctx context.Context, keyvals ...interface{}) error {
	return level.Debug(contextual.GetLogger(ctx)).Log(keyvals...)
}
