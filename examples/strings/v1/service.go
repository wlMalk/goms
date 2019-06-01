package strings

import (
	"context"
)

// Service is the interface describing the strings service
// @http-URI-prefix(/strings)
// @gen(middleware, logging, grpc, http, recover, main, error-logging)
type Service interface {
	// @http-method(POST)
	// @http-URI(/uppercase)
	// @http-abs-URI(/strings/getUpperCase)
	// @params([strs], (http-origin(BODY) ) )
	// @logs-ignore(ans, err)
	// @alias(strs,strings) @alias(ans,answer)
	// @validate
	Uppercase(ctx context.Context, strs ...string) (ans []string, err error)
	Lowercase(ctx context.Context, strs ...string) (ans []string, err error)
}
