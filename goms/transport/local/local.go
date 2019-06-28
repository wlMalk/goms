package local

import (
	"context"
	"strings"

	"github.com/wlMalk/goms/goms/correlation"
	"github.com/wlMalk/goms/goms/log/contextual"
	"github.com/wlMalk/goms/goms/request"
	"github.com/wlMalk/goms/goms/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Client struct {
	endpoint  endpoint.Endpoint
	before    []RequestFunc
	finalizer []FinalizerFunc
}

func New(
	endpoint endpoint.Endpoint,
	options ...Option,
) *Client {
	c := &Client{
		endpoint: endpoint,
		before:   []RequestFunc{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type Option func(*Client)

func Before(before ...RequestFunc) Option {
	return func(c *Client) { c.before = append(c.before, before...) }
}

func Finalizer(f ...FinalizerFunc) Option {
	return func(s *Client) { s.finalizer = append(s.finalizer, f...) }
}

func (c Client) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		oCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if c.finalizer != nil {
			defer func() {
				for _, f := range c.finalizer {
					f(ctx, err)
				}
			}()
		}

		for _, f := range c.before {
			oCtx = f(ctx, oCtx, request)
		}

		return c.endpoint(oCtx, request)
	}
}

type RequestFunc func(ctx context.Context, oCtx context.Context, req interface{}) context.Context

type FinalizerFunc func(ctx context.Context, err error)

func LoggerInjector(logger log.Logger) RequestFunc {
	return func(ctx context.Context, oCtx context.Context, req interface{}) context.Context {
		requestID := request.GetRequestID(oCtx)
		logger = contextual.ContextualLogger(oCtx, logger, "request-id", requestID)
		return contextual.SetLogger(oCtx, logger)
	}
}

func RequestIDCreator() RequestFunc {
	return func(ctx context.Context, oCtx context.Context, req interface{}) context.Context {
		requestID := request.NewRequestID()
		return request.SetRequestID(oCtx, requestID)
	}
}

func MethodInjector(ser string, met string) RequestFunc {
	return func(ctx context.Context, oCtx context.Context, req interface{}) context.Context {
		method := service.NewMethod(ser, met)
		return service.SetMethod(oCtx, method)
	}
}

func CorrelationIDInjector() RequestFunc {
	return func(ctx context.Context, oCtx context.Context, req interface{}) context.Context {
		correlationID := correlation.GetCorrelationID(ctx)
		if len(strings.TrimSpace(correlationID)) > 0 {
			oCtx = correlation.SetCorrelationID(oCtx, correlationID)
		}
		return oCtx
	}
}

func RequestIDInjector() RequestFunc {
	return func(ctx context.Context, oCtx context.Context, req interface{}) context.Context {
		requestID := request.GetRequestID(ctx)
		if len(strings.TrimSpace(requestID)) > 0 {
			oCtx = request.SetCallerRequestID(oCtx, requestID)
		}
		return oCtx
	}
}
