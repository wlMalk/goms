package request

import (
	"context"
)

type contextKeyType string

const contextRequestIDKey contextKeyType = "request-id"

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
