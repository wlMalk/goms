package http

import (
	"context"
	"encoding/json"
	"net/http"
)

func HTTPResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type Router interface {
	Method(method string, uri string, handler http.Handler)
	Params(req *http.Request) Params
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Params interface {
	Get(name string) string
}
