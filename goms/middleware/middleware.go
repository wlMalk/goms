package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/wlMalk/goms/goms/log"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

func LatencyMiddleware(latency metrics.Histogram) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			defer func(begin time.Time) {
				latency.With("service", "service", "method", "method").Observe(time.Since(begin).Seconds())
			}(time.Now())

			return e(ctx, req)
		}
	}
}

