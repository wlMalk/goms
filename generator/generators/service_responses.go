package generators

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceResponseStruct(file file.File, service types.Service, method types.Method) error {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("type %sResponse struct {", methodName)
	for _, res := range method.Results {
		file.Pf("%s %s", strings.ToUpperFirst(res.Name), res.Type.GoType())
	}
	file.Pf("}")
	return nil
}
