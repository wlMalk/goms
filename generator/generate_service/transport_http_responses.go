package generate_service

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func HTTPResponse(file file.File, service types.Service, method types.Method) {
	if len(method.Results) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("type %sResponse struct {", methodName)
	HTTPResponseFields(file, method.Results)
	file.Pf("}")
	file.Pf("")
}

func HTTPResponseFields(file file.File, fields []*types.Field) {
	for _, field := range fields {
		fieldName := strings.ToUpperFirst(field.Name)
		fieldSpecialName := helpers.GetName(strings.ToLowerFirst(field.Name), field.Alias)
		file.Pf("%s %s `json:\"%s\"`", fieldName, field.Type.GoType(), fieldSpecialName)
	}
}

func HTTPResponseNewFunc(file file.File, service types.Service, method types.Method) {
	if len(method.Results) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", service.ImportPath, "/pkg/service/responses")
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

func HTTPResponseNewHTTPFunc(file file.File, service types.Service, method types.Method) {
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

func HTTPResponseToResponseFunc(file file.File, service types.Service, method types.Method) {
	if len(method.Results) == 0 {
		return
	}
	file.AddImport("", service.ImportPath, "/pkg/service/responses")
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
