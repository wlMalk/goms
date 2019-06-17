package generators

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceStructType(file file.File, service types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	serviceName := strings.ToUpperFirst(service.Name)
	lowerServiceName := strings.ToLowerFirst(service.Name)
	file.Pf("type endpointsHandler struct {")
	for _, method := range service.Methods {
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s endpoint.Endpoint", lowerMethodName)
	}
	file.Pf("}")
	file.Pf("")
	file.Pf("type %s struct {", lowerServiceName)
	file.Pf("endpoints *endpointsHandler")
	file.Pf("handler handlers.EndpointHandler")
	file.Pf("}")
	file.Pf("")
	file.Pf("type %s struct {", serviceName)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%s endpoint.Endpoint", methodName)
	}
	file.Pf("}")
	file.Pf("")
}

func ServiceStructTypeNewFunc(file file.File, service types.Service) {
	file.AddImport("", "context")
	file.AddImport("", "github.com/wlMalk/goms/goms/errors")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	serviceName := strings.ToUpperFirst(service.Name)
	lowerServiceName := strings.ToLowerFirst(service.Name)
	file.Pf("func Endpoints(h interface{}, middlewareGetter interface{}) %s {", serviceName)
	file.Pf("handler := &endpointsHandler{")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: endpoint.Endpoint(func(ctx context.Context, req interface{}) (res interface{}, err error) {", lowerMethodName)
		file.Pf("return nil, errors.MethodNotImplemented(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}),")
	}
	file.Pf("}")
	file.Pf("s := &%s {", lowerServiceName)
	file.P("endpoints: handler,")
	file.P("handler:   handler,")
	file.Pf("}")
	file.Pf("")
	ServiceMethodsRegisteration(file, service)
	file.Pf("}")
	file.Pf("")
}

func ServiceMethodsRegisteration(file file.File, service types.Service) {
	for _, method := range service.Methods {
		TypeSwitchForMethodHandler(file, service, method)
		if method.Options.Generate.Middleware {
			MiddlewareCheckerForEndpoint(file, service, method)
		}
	}
	if service.Options.Generate.Middleware {
		MiddlewareCheckerForService(file, service)
	}
	EndpointsPacker(file, service)
}

func TypeSwitchForMethodHandler(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers/converters")
	file.Pf("switch t := h.(type) {")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("case handlers.%sHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.%sRequestResponseHandlerToEndpoint(converters.%sHandlerTo%sRequestResponseHandler(handlers.%sHandlerFunc(t.%s)))", lowerMethodName, methodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestResponseHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.%sRequestResponseHandlerToEndpoint(handlers.%sRequestResponseHandlerFunc(t.%s))", lowerMethodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sEndpointHandler:", methodName)
	file.Pf("s.endpoints.%s = t.%s", lowerMethodName, methodName)
	file.Pf("}")
	file.Pf("")
}

func MiddlewareCheckerForEndpoint(file file.File, service types.Service, method types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("if t, ok := middlewareGetter.(interface{ %sMiddleware(e endpoint.Endpoint) endpoint.Endpoint }); ok {", methodName)
	file.Pf("s.endpoints.%s = t.%sMiddleware(s.endpoints.%s)", lowerMethodName, methodName, lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func MiddlewareCheckerForService(file file.File, service types.Service) {
	file.Pf("if t, ok := middlewareGetter.(interface{ Middleware(h handlers.EndpointHandler) handlers.EndpointHandler }); ok {")
	file.Pf("s.handler = t.Middleware(s.endpoints)")
	file.Pf("}")
	file.Pf("")
}

func OuterMiddlewareCheckerForEndpoint(file file.File, service types.Service, method types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("if t, ok := h.(interface {")
	file.Pf("Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint", methodName)
	file.Pf("}); ok {")
	file.Pf("%s = t.Outer%sMiddleware(%s)", lowerMethodName, methodName, lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func ServiceStructMethodHandler(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("func (s *endpointsHandler) %s(ctx context.Context, req interface{}) (res interface{}, err error) {", methodName)
	file.Pf("return s.%s(ctx, req)", lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func EndpointsPacker(file file.File, service types.Service) {
	file.Pf("var (")
	serviceName := strings.ToUpperFirst(service.Name)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s endpoint.Endpoint = s.handler.%s", lowerMethodName, methodName)
	}
	file.Pf(")")
	file.Pf("")
	for _, method := range helpers.GetMethodsWithMiddlewareEnabled(service) {
		OuterMiddlewareCheckerForEndpoint(file, service, method)
	}
	file.Pf("endpoints := %s{", serviceName)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: %s,", methodName, lowerMethodName)
	}
	file.Pf("}")
	file.Pf("")
	file.Pf("return endpoints")
}
