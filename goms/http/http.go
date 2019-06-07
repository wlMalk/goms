package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func HTTPResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type Router interface {
	Method(method string, uri string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Params interface {
	Get(name string) string
}

type contextKey int

const ContextParamsKey contextKey = iota

func SetParams(ctx context.Context, params Params) context.Context {
	return context.WithValue(ctx, ContextParamsKey, params)
}

func GetParams(ctx context.Context) Params {
	return ctx.Value(ContextParamsKey).(Params)
}

