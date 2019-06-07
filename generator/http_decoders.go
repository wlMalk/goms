package generator

import (
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPDecodersFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	for _, method := range service.Methods {
		generateHTTPRequestDecoder(file, method)
		generateHTTPResponseDecoder(file, method)
	}
	return file
}

func generateHTTPRequestDecoder(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sRequest(ctx context.Context, req *http.Request) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "service/transport/http/requests")
		file.Pf("r, err := requests.%sFromHTTP(req)", methodName)
		file.Pf("if err!=nil{")
		file.Pf("return nil, err")
		file.Pf("}")
		file.Pf("return r.Request(), err")
	} else {
		file.Pf("return nil, nil")
	}
	file.Pf("}")
	file.Pf("")
}

func generateHTTPResponseDecoder(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sResponse(ctx context.Context, res *http.Response) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "service/transport/http/responses")
		file.Pf("resp, err := responses.%sFromHTTP(res)", methodName)
		file.Pf("if err!=nil{")
		file.Pf("return nil, err")
		file.Pf("}")
		file.Pf("return resp.Response(), err")
	} else {
		file.Pf("return nil, nil")
	}
	file.Pf("}")
	file.Pf("")
}
