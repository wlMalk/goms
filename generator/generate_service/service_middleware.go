package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateServiceMiddlewareFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateServiceMiddlewareTypes(file, service)
	generateServiceMiddlewareChainFunc(file, service)
	generateServiceRequestMiddlewareChainFunc(file, service)
	generateServiceRequestResponseMiddlewareChainFunc(file, service)
	generateServiceEndpointMiddlewareStruct(file, service)
	generateServiceEndpointMiddlewareChainFunc(file, service)
	generateServiceEndpointMiddlewareNewFuncs(file, service)
	for _, method := range service.Methods {
		generateServiceEndpointMiddlewareMethod(file, method)
	}
	return file
}

func generateServiceMiddlewareTypes(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.P("type (")
	file.P("Middleware                func(handlers.Handler) handlers.Handler")
	file.P("RequestMiddleware         func(handlers.RequestHandler) handlers.RequestHandler")
	file.P("RequestResponseMiddleware func(handlers.RequestResponseHandler) handlers.RequestResponseHandler")
	file.P(")")
	file.P("")
}

func generateServiceMiddlewareChainFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
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

func generateServiceRequestMiddlewareChainFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.P("func ChainRequest(outer RequestMiddleware, others ...RequestMiddleware) RequestMiddleware {")
	file.P("return func(next handlers.RequestHandler) handlers.RequestHandler {")
	file.P("for i := len(others) - 1; i >= 0; i-- {")
	file.P("next = others[i](next)")
	file.P("}")
	file.P("return outer(next)")
	file.P("}")
	file.P("}")
	file.P("")
}

func generateServiceRequestResponseMiddlewareChainFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
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

func generateServiceEndpointMiddlewareStruct(file *files.GoFile, service *types.Service) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.Pf("type endpointMiddleware struct {")
	for _, method := range service.Methods {
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s endpoint.Endpoint", lowerMethodName)
	}
	file.Pf("}")
	file.Pf("")
}

func generateServiceEndpointMiddlewareChainFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.Pf("func chainEndpointMiddleware(mw ...endpoint.Middleware) endpoint.Middleware {")
	file.Pf("if len(mw) == 1 {")
	file.Pf("return mw[0]")
	file.Pf("} else if len(mw) > 1 {")
	file.Pf("return endpoint.Chain(mw[0], mw[1:]...)")
	file.Pf("}")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateServiceEndpointMiddlewareNewFuncs(file *files.GoFile, service *types.Service) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.Pf("func EndpointMiddleware(h handlers.EndpointHandler, mw ...endpoint.Middleware) handlers.EndpointHandler {")
	file.Pf("return EndpointMiddlewareSpecial(h, func(method string) bool {")
	file.Pf("return true")
	file.Pf("}, mw...)")
	file.Pf("}")
	file.Pf("")
	file.Pf("func EndpointMiddlewareSpecial(h handlers.EndpointHandler, f func(method string) bool, mw ...endpoint.Middleware) handlers.EndpointHandler {")
	file.Pf("if len(mw) == 0 {")
	file.Pf("return h")
	file.Pf("}")
	file.Pf("")
	file.Pf("fun := func(name string, e endpoint.Endpoint) endpoint.Endpoint {")
	file.Pf("if !f(name) {")
	file.Pf("return e")
	file.Pf("}")
	file.Pf("return chainEndpointMiddleware(mw...)(e)")
	file.Pf("}")
	file.Pf("")
	file.Pf("return &endpointMiddleware{")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: fun(\"%s\", h.%s),", lowerMethodName, methodName, methodName)
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}

func generateServiceEndpointMiddlewareMethod(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("func (mw endpointMiddleware) %s(ctx context.Context, req interface{}) (res interface{}, err error) {", methodName)
	file.Pf("return mw.%s(ctx, req)", lowerMethodName)
	file.Pf("}")
	file.Pf("")
}
