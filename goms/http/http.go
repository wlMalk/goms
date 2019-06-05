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
	PathParams(req *http.Request) PathParams
}

type PathParams interface {
	Get(name string) string
}
