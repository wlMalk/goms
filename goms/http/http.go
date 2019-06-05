package http

import (
	"context"
	"encoding/json"
	"net/http"
)

type Router interface {
	Method(method string, uri string, handler http.Handler)
	PathParams(req *http.Request) PathParams
}

type PathParams interface {
	Get(name string) string
}
