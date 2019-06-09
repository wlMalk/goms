package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateHTTPResponsesFile(base string, path string, name string, methods []*types.Method) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range methods {
		generateHTTPResponse(file, method)
		generateHTTPResponseNewFunc(file, method)
		generateHTTPResponseNewHTTPFunc(file, method)
		generateHTTPResponseToResponseFunc(file, method)
	}
	return file
}

func generateHTTPResponse(file *files.GoFile, method *types.Method) {
	if len(method.Results) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("type %sResponse struct {", methodName)
	generateHTTPResponseFields(file, method.Results)
	file.Pf("}")
	file.Pf("")
}

func generateHTTPResponseFields(file *files.GoFile, fields []*types.Field) {
	for _, field := range fields {
		fieldName := strings.ToUpperFirst(field.Name)
		fieldSpecialName := helpers.GetName(field.Name, field.Alias)
		file.Pf("%s %s `json:\"%s\"`", fieldName, field.Type.GoType(), fieldSpecialName)
	}
}

func generateHTTPResponseNewFunc(file *files.GoFile, method *types.Method) {
	if len(method.Results) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "service/responses")
	file.Pf("func %s(res *responses.%sResponse) *%sResponse {", methodName, methodName, methodName)
	file.Pf("resp := &%sResponse{}", methodName)
	for _, res := range method.Results {
		resName := strings.ToUpperFirst(res.Name)
		file.Pf("resp.%s = res.%s", resName, resName)
	}
	file.Pf("return resp")
	file.Pf("}")
	file.Pf("")
}

func generateHTTPResponseNewHTTPFunc(file *files.GoFile, method *types.Method) {
	if len(method.Results) == 0 {
		return
	}
	file.AddImport("", "net/http")
	file.AddImport("", "encoding/json")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func %sFromHTTP(res *http.Response) (*%sResponse, error) {", methodName, methodName)
	file.Pf("resp := &%sResponse{}", methodName)
	file.Pf("err := json.NewDecoder(res.Body).Decode(resp)")
	file.Pf("if err!=nil{")
	file.Pf("return nil, err")
	file.Pf("}")
	file.Pf("return resp, nil")
	file.Pf("}")
	file.Pf("")
}

func generateHTTPResponseToResponseFunc(file *files.GoFile, method *types.Method) {
	if len(method.Results) == 0 {
		return
	}
	file.AddImport("", method.Service.ImportPath, "service/responses")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func (r *%sResponse) Response() *responses.%sResponse {", methodName, methodName)
	file.Pf("resp := &responses.%sResponse{}", methodName)
	for _, res := range method.Results {
		resName := strings.ToUpperFirst(res.Name)
		file.Pf("resp.%s = r.%s", resName, resName)
	}
	file.Pf("return resp")
	file.Pf("}")
	file.Pf("")
}
