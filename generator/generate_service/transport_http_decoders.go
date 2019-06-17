package generate_service

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func HTTPRequestDecoder(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sRequest(ctx context.Context, req *http.Request) (interface{}, error) {", methodName)
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/transport/http/requests")
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

func HTTPResponseDecoder(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sResponse(ctx context.Context, res *http.Response) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/transport/http/responses")
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
