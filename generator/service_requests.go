package generator

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateRequestsFile(base string, path string, name string, methods []*types.Method) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	for _, m := range methods {
		var fields []*types.Field
		for _, arg := range m.Arguments {
			f := &types.Field{}
			f.Name = arg.Name
			f.Type = arg.Type
			f.Alias = arg.Alias
			fields = append(fields, f)
		}
		if len(fields) == 0 {
			continue
		}
		generateExportedStruct(file, m.Name+"Request", fields)
		generateMethodRequestNewFunc(file, m)
	}
	return file
}

func generateMethodRequestNewFunc(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	args := getMethodArguments(method.Arguments)
	file.Pf("func %s(%s)*%sRequest{", methodName, strs.Join(args, ", "), methodName)
	file.Pf("return &%sRequest{", methodName)
	for _, arg := range method.Arguments {
		file.Pf("%s: %s,", strings.ToUpperFirst(arg.Name), strings.ToLowerFirst(arg.Name))
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}
