package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/wlMalk/goms/goms/correlation"
	"github.com/wlMalk/goms/goms/log"
	"github.com/wlMalk/goms/goms/request"
	"github.com/wlMalk/goms/goms/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

func LatencyMiddleware(latency metrics.Histogram) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			defer func(begin time.Time) {
				latency.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())

			return e(ctx, req)
		}
	}
}

func CounterMiddleware(counter metrics.Counter) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			counter.Add(1)
			return e(ctx, req)
		}
	}
}

func FrequencyMiddleware(frequency metrics.Gauge) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			frequency.Add(1)
			res, err = e(ctx, req)
			frequency.Add(-1)
			return
		}
	}
}

func InstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, frequency metrics.Gauge) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			defer func(begin time.Time) {
				latency.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
				frequency.Add(-1)
			}(time.Now())
			frequency.Add(1)
			counter.Add(1)
			res, err = e(ctx, req)
			return
		}
	}
}

func LoggingMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			method := service.GetMethod(ctx)
			correlationID := correlation.GetCorrelationID(ctx)
			callerRequestID := request.GetCallerRequestID(ctx)
			defer func(begin time.Time) {
				if err != nil {
					log.Log(ctx, "service", method.Service.Name, "method", method.Name, "correlation_id", correlationID, "caller_request_id", callerRequestID, "transport_error", err, "took", time.Since(begin))
				} else {
					log.Log(ctx, "service", method.Service.Name, "method", method.Name, "correlation_id", correlationID, "caller_request_id", callerRequestID, "took", time.Since(begin))
				}
			}(time.Now())
			return e(ctx, req)

		}
	}
}

func RecoveringMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
				}
			}()
			return e(ctx, req)
		}
	}
}
