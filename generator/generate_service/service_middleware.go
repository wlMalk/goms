package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateServiceMiddlewareFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateServiceMiddlewareTypes(file, service)
	generateServiceMiddlewareChainFunc(file, service)
	generateServiceRequestResponseMiddlewareChainFunc(file, service)
	generateServiceApplyMiddlewareFunc(file, service)
	generateServiceApplyMiddlewareSpecialFunc(file, service)
	generateServiceApplyMiddlewareConditionalFunc(file, service)
	return file
}

func generateServiceMiddlewareTypes(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.P("type (")
	file.P("Middleware                func(handlers.Handler) handlers.Handler")
	file.P("RequestResponseMiddleware func(handlers.RequestResponseHandler) handlers.RequestResponseHandler")
	file.P(")")
	file.P("")
}

func generateServiceMiddlewareChainFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.P("func Chain(outer Middleware, others ...Middleware) Middleware {")
	file.P("return func(next handlers.Handler) handlers.Handler {")
	file.P("for i := len(others) - 1; i >= 0; i-- {")
	file.P("next = others[i](next)")
	file.P("}")
	file.P("return outer(next)")
	file.P("}")
	file.P("}")
	file.P("")
}

func generateServiceRequestResponseMiddlewareChainFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.P("func ChainRequestResponse(outer RequestResponseMiddleware, others ...RequestResponseMiddleware) RequestResponseMiddleware {")
	file.P("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
	file.P("for i := len(others) - 1; i >= 0; i-- {")
	file.P("next = others[i](next)")
	file.P("}")
	file.P("return outer(next)")
	file.P("}")
	file.P("}")
	file.P("")
}

func generateServiceApplyMiddlewareFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.AddImport("", service.ImportPath, "/pkg/transport")
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func ApplyMiddleware(endpoints transport.%s, mw ...endpoint.Middleware) transport.%s {", serviceName, serviceName)
	file.Pf("return ApplyMiddlewareConditional(endpoints, func(method string) bool {")
	file.Pf("return true")
	file.Pf("}, mw...)")
	file.Pf("}")
	file.Pf("")
}

func generateServiceApplyMiddlewareSpecialFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.AddImport("", service.ImportPath, "/pkg/transport")
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func ApplyMiddlewareSpecial(endpoints transport.%s, middlewareFunc func(method string) (mw []endpoint.Middleware)) transport.%s {", serviceName, serviceName)
	file.Pf("if middlewareFunc == nil {")
	file.Pf("return endpoints")
	file.Pf("}")
	file.Pf("")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("endpoints.%s = goms_endpoint.Chain(middlewareFunc(\"%s\")...)(endpoints.%s)", methodName, helpers.GetName(methodName, method.Alias), methodName)
	}
	file.Pf("")
	file.Pf("return endpoints")
	file.Pf("}")
	file.Pf("")
}
func generateServiceApplyMiddlewareConditionalFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.AddImport("", service.ImportPath, "/pkg/transport")
	file.AddImport("goms_endpoint", "github.com/wlMalk/goms/goms/endpoint")
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func ApplyMiddlewareConditional(endpoints transport.%s, f func(method string) bool, mw ...endpoint.Middleware) transport.%s {", serviceName, serviceName)
	file.Pf("if len(mw) == 0 {")
	file.Pf("return endpoints")
	file.Pf("}")
	file.Pf("")
	file.Pf("fun := func(method string, e endpoint.Endpoint) endpoint.Endpoint {")
	file.Pf("if !f(method) {")
	file.Pf("return e")
	file.Pf("}")
	file.Pf("return goms_endpoint.Chain(mw...)(e)")
	file.Pf("}")
	file.Pf("")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("endpoints.%s = fun(\"%s\", endpoints.%s)", methodName, helpers.GetName(methodName, method.Alias), methodName)
	}
	file.Pf("return endpoints")
	file.Pf("}")
	file.Pf("")
}
