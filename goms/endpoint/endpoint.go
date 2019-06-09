package endpoint

import (
	"github.com/go-kit/kit/endpoint"
)

func Chain(mw ...endpoint.Middleware) endpoint.Middleware {
	if len(mw) == 1 {
		return mw[0]
	} else if len(mw) > 1 {
		return endpoint.Chain(mw[0], mw[1:]...)
	}
	return nil
}
