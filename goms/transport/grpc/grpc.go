package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/wlMalk/goms/goms/correlation"
	"github.com/wlMalk/goms/goms/log/contextual"
	"github.com/wlMalk/goms/goms/request"
	"github.com/wlMalk/goms/goms/service"

	"github.com/go-kit/kit/log"
	kit_grpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	Listener net.Listener
	Server   *grpc.Server
}

func NewServer(listener net.Listener, opts ...grpc.ServerOption) *Server {
	return &Server{
		Listener: listener,
		Server:   grpc.NewServer(opts...),
	}
}

func LoggerInjector(logger log.Logger) kit_grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		requestID := request.GetRequestID(ctx)
		logger = contextual.ContextualLogger(ctx, logger, "request-id", requestID)
		return contextual.SetLogger(ctx, logger)
	}
}

func RequestIDCreator() kit_grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		requestID := request.NewRequestID()
		return request.SetRequestID(ctx, requestID)
	}
}

func CorrelationIDExtractor() kit_grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		s := md.Get("Correlation-ID")
		var correlationID string
		if len(s) > 0 {
			correlationID = s[0]
		} else {
			correlationID = correlation.NewCorrelationID()
		}
		return correlation.SetCorrelationID(ctx, correlationID)
	}
}

func RequestIDExtractor() kit_grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		s := md.Get("Caller-Request-ID")
		var requestID string
		if len(s) > 0 {
			requestID = s[0]
		} else {
			return ctx
		}
		return request.SetCallerRequestID(ctx, requestID)
	}
}

func CorrelationIDInjector() kit_grpc.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		correlationID := correlation.GetCorrelationID(ctx)
		if len(strings.TrimSpace(correlationID)) > 0 {
			md.Set("Correlation-ID", correlationID)
		}
		return ctx
	}
}

func RequestIDInjector() kit_grpc.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		requestID := request.GetRequestID(ctx)
		if len(strings.TrimSpace(requestID)) > 0 {
			md.Set("Caller-Request-ID", requestID)
		}
		return ctx
	}
}

func MethodInjector(ser string, met string) kit_grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		method := service.NewMethod(ser, met)
		return service.SetMethod(ctx, method)
	}
}
