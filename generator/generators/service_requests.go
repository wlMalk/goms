package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func MethodRequestNewFunc(file file.File, service types.Service, method types.Method) {
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
}
