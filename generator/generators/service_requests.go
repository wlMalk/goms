package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceRequestStruct(file file.File, service types.Service, method types.Method) error {
	helpers.AddTypesImports(file, service)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("type %sRequest struct {", methodName)
	for _, arg := range method.Arguments {
		file.Pf("%s %s", strings.ToUpperFirst(arg.Name), arg.Type.GoType())
	}
	file.Pf("}")
	return nil
}

func ServiceRequestNewFunc(file file.File, service types.Service, method types.Method) error {
	helpers.AddTypesImports(file, service)
	methodName := strings.ToUpperFirst(method.Name)
	args := helpers.GetMethodArguments(method.Arguments)
	file.Pf("func %s(%s)*%sRequest{", methodName, strs.Join(args, ", "), methodName)
	file.Pf("return &%sRequest{", methodName)
	for _, arg := range method.Arguments {
		file.Pf("%s: %s,", strings.ToUpperFirst(arg.Name), strings.ToLowerFirst(arg.Name))
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
	return nil
}
