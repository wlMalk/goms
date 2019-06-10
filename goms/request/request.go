package request

import (
	"context"

	"github.com/rs/xid"
)

type contextKeyType string

const (
	contextRequestIDKey       contextKeyType = "request-id"
	contextCallerRequestIDKey contextKeyType = "caller-request-id"
)

func NewRequestID() string {
	return xid.New().String()
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextRequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	requestID := ctx.Value(contextRequestIDKey)
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}

func SetCallerRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextCallerRequestIDKey, requestID)
}

func GetCallerRequestID(ctx context.Context) string {
	requestID := ctx.Value(contextCallerRequestIDKey)
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}
