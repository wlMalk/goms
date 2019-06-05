package generator

import (
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPResponsesFile(base string, path string, name string, methods []*types.Method) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	for _, method := range methods {
		generateHTTPResponse(file, method)
	}
	return file
}

func generateHTTPResponse(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("type %sResponse struct {", methodName)
	generateHTTPResponseFields(file, method.Results)
	file.Pf("}")
	file.Pf("")
}
func generateHTTPResponseFields(file *GoFile, fields []*types.Field) {
	for _, field := range fields {
		fieldName := strings.ToUpperFirst(field.Name)
		lowerFieldName := strings.ToLowerFirst(field.Name)
		file.Pf("%s %s `json:\"%s\"`", fieldName, field.Type.GoType(), lowerFieldName)
	}
}
