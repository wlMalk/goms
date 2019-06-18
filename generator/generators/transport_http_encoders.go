package generators

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func HTTPRequestEncoder(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	file.AddImport("", "path")
	serviceName := strings.ToUpperFirst(service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	methodHTTPMethod := method.Options.HTTP.Method
	methodURI := getMethodURI(service, method)
	file.Pf("func Encode%sRequest(ctx context.Context, r *http.Request, request interface{}) error {", methodName)
	file.Pf("r.Method = \"%s\"", methodHTTPMethod)
	file.Pf("r.URL.Path = path.Join(r.URL.Path, \"%s\")", methodURI)
	if len(method.Arguments) > 0 {
		file.AddImport("http_requests", service.ImportPath, "/pkg/transport/http/requests")
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if request == nil {")
		file.Pf("return errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return http_requests.%s(request.(*requests.%sRequest)).ToHTTP(r)", methodName, methodName)
	} else {
		file.Pf("return nil")
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func HTTPResponseEncoder(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", "net/http")
	serviceName := strings.ToUpperFirst(service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
		file.AddImport("http_responses", service.ImportPath, "/pkg/transport/http/responses")
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if response == nil {")
		file.Pf("return errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return goms_http.HTTPResponseEncoder(ctx, w, http_responses.%s(response.(*responses.%sResponse)))", methodName, methodName)
	} else {
		file.Pf("return nil")
	}
	file.Pf("}")
	file.Pf("")
	return nil
}
