package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateHTTPEncodersFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range helpers.GetMethodsWithHTTPEnabled(service) {
		generateHTTPRequestEncoder(file, method)
		generateHTTPResponseEncoder(file, method)
	}
	return file
}

func generateHTTPRequestEncoder(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sRequest(ctx context.Context, r *http.Request, request interface{}) error {", methodName)
	if len(method.Arguments) > 0 {
		file.AddImport("http_requests", method.Service.ImportPath, "/service/transport/http/requests")
		file.AddImport("", method.Service.ImportPath, "service/requests")
		file.Pf("return http_requests.%s(request.(*requests.%sRequest)).ToHTTP(r)", methodName, methodName)
	} else {
		file.Pf("return nil")
	}
	file.Pf("}")
	file.Pf("")
}

func generateHTTPResponseEncoder(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
		file.AddImport("http_responses", method.Service.ImportPath, "/service/transport/http/responses")
		file.AddImport("", method.Service.ImportPath, "service/responses")
		file.Pf("return goms_http.HTTPResponseEncoder(ctx, w, http_responses.%s(response.(*responses.%sResponse)))", methodName, methodName)
	} else {
		file.Pf("return nil")
	}
	file.Pf("}")
	file.Pf("")
}
