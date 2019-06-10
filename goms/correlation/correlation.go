package correlation

import (
	"context"

	"github.com/rs/xid"
)

type contextKeyType string

const contextCorrelationIDKey contextKeyType = "correlation-id"

func NewCorrelationID() string {
	return xid.New().String()
}

func SetCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, contextCorrelationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	correlationID := ctx.Value(contextCorrelationIDKey)
	if correlationID == nil {
		return ""
	}
	return correlationID.(string)
}
