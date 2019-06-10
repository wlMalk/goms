package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wlMalk/goms/goms/correlation"
	"github.com/wlMalk/goms/goms/log/contextual"
	"github.com/wlMalk/goms/goms/request"
	"github.com/wlMalk/goms/goms/service"

	"github.com/go-kit/kit/log"
	kit_http "github.com/go-kit/kit/transport/http"
)

func HTTPResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type Server struct {
	router Router
	http.Server
}

func (s *Server) RegisterMethod(method string, uri string, handler http.Handler) {
	s.router.Method(method, uri, handler)
}

func NewServer(router Router) *Server {
	return &Server{router: router}
}

type Router interface {
	Method(method string, uri string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Params interface {
	Get(name string) string
}

type contextKeyType string

const contextParamsKey contextKeyType = "params"

func SetParams(ctx context.Context, params Params) context.Context {
	return context.WithValue(ctx, contextParamsKey, params)
}

func GetParams(ctx context.Context) Params {
	params := ctx.Value(contextParamsKey)
	if params == nil {
		return nil
	}
	return params.(Params)
}

func FormatURI(uri string, pairs ...string) string {
	if len(pairs)%2 != 0 {
		return ""
	}
	uri += "/"
	for i := 0; i < len(pairs); i += 2 {
		strings.Replace(uri, "/:"+pairs[i]+"/", "/"+pairs[i+1]+"/", -1)
		strings.Replace(uri, "/*"+pairs[i]+"/", "/"+pairs[i+1]+"/", -1)
	}
	return strings.TrimSuffix(uri, "/")
}

func LoggerInjector(logger log.Logger) kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		requestID := request.GetRequestID(ctx)
		logger = contextual.ContextualLogger(ctx, logger, "request-id", requestID)
		return contextual.SetLogger(ctx, logger)
	}
}

func RequestIDCreator() kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		requestID := request.NewRequestID()
		return request.SetRequestID(ctx, requestID)
	}
}

func CorrelationIDExtractor() kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		correlationID := r.Header.Get("X-Correlation-ID")
		if len(strings.TrimSpace(correlationID)) == 0 {
			correlationID = correlation.NewCorrelationID()
		}
		return correlation.SetCorrelationID(ctx, correlationID)
	}
}

func RequestIDExtractor() kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		requestID := r.Header.Get("X-Caller-Request-ID")
		if len(strings.TrimSpace(requestID)) == 0 {
			return ctx
		}
		return request.SetCallerRequestID(ctx, requestID)
	}
}

func CorrelationIDInjector() kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		correlationID := correlation.GetCorrelationID(ctx)
		if len(strings.TrimSpace(correlationID)) > 0 {
			r.Header.Set("X-Correlation-ID", correlationID)
		}
		return ctx
	}
}

func RequestIDInjector() kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		requestID := request.GetRequestID(ctx)
		if len(strings.TrimSpace(requestID)) > 0 {
			r.Header.Set("X-Caller-Request-ID", requestID)
		}
		return ctx
	}
}

func MethodInjector(ser string, met string) kit_http.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		method := service.NewMethod(ser, met)
		return service.SetMethod(ctx, method)
	}
}
