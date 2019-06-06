package generator

import (
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPEncodersFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	for _, method := range service.Methods {
		generateHTTPRequestEncoder(file, method)
		generateHTTPResponseEncoder(file, method)
	}
	return file
}

func generateHTTPRequestEncoder(file *GoFile, method *types.Method) {

}

func generateHTTPResponseEncoder(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/http")
		file.AddImport("http_responses", method.Service.ImportPath, "/service/transport/http/responses")
		file.AddImport("", method.Service.ImportPath, "service/responses")
		file.Pf("return goms_http.HTTPResponseEncoder(ctx, w, http_responses.%s(response.(*responses.%sResponse)))", methodName, methodName)
	} else {
		file.Pf("return nil")
	}
	file.Pf("}")
	file.Pf("")
}
