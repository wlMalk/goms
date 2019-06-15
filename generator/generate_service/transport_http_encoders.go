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
	file.AddImport("", "path")
	serviceName := strings.ToUpperFirst(method.Service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	methodHTTPMethod := method.Options.HTTP.Method
	methodURI := getMethodURI(method)
	file.Pf("func Encode%sRequest(ctx context.Context, r *http.Request, request interface{}) error {", methodName)
	file.Pf("r.Method = \"%s\"", methodHTTPMethod)
	file.Pf("r.URL.Path = path.Join(r.URL.Path, \"%s\")", methodURI)
	if len(method.Arguments) > 0 {
		file.AddImport("http_requests", method.Service.ImportPath, "/pkg/transport/http/requests")
		file.AddImport("", method.Service.ImportPath, "/pkg/service/requests")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if request == nil {")
		file.Pf("return errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, method.Service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
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
	serviceName := strings.ToUpperFirst(method.Service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
		file.AddImport("http_responses", method.Service.ImportPath, "/pkg/transport/http/responses")
		file.AddImport("", method.Service.ImportPath, "/pkg/service/responses")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if response == nil {")
		file.Pf("return errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, method.Service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return goms_http.HTTPResponseEncoder(ctx, w, http_responses.%s(response.(*responses.%sResponse)))", methodName, methodName)
	} else {
		file.Pf("return nil")
	}
	file.Pf("}")
	file.Pf("")
}
