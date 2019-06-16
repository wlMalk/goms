package tags

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/parser/types"
)

func ParamHTTPOriginTag(arg *types.Argument, tag string) error {
	origin := strs.ToUpper(tag)
	if origin != "BODY" && origin != "HEADER" && origin != "QUERY" && origin != "PATH" {
		return fmt.Errorf("invalid http-origin value '%s'", tag)
	}
	arg.Options.HTTP.Origin = origin
	return nil
}
