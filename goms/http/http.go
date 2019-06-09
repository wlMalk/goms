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
