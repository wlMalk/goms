package strings

import (
	"context"
)

// Service is the interface describing the strings service
// @http-URI-prefix(/strings)
// @name(com.goms.strings.v0.5.0)
// @generate-all(grpc)
type Strings_v0_5 interface {
	/*
		@name(upper)
		@http-method(GET)
		@http-abs-URI(/strings/uc)
		@params([strs], (@http-origin(QUERY)))
	*/
	/*
		@logs-ignore(
			strs,
			ans,
		)
	*/
	// @alias(strs,strings) @alias(ans,answer)
	Uppercase(ctx context.Context, strs ...string) (ans []string, err error)
}
