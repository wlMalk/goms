package httprouter

import (
	"net/http"

	goms_http "github.com/wlMalk/goms/goms/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	r *httprouter.Router
}

type Params struct {
	params httprouter.Params
}

func (p Params) Get(name string) string {
	return p.params.ByName(name)
}

func New(r *httprouter.Router) *Router {
	return &Router{
		r: r,
	}
}

func (r *Router) Method(method string, uri string, handler http.Handler) {
	r.r.Handle(method, uri, httprouter.Handle(func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		req = req.WithContext(goms_http.SetParams(req.Context(), &Params{params: params}))
		handler.ServeHTTP(w, req)
	}))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}
