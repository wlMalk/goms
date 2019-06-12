package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateServiceTransportEndpointsFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateServiceStructType(file, service)
	generateServiceStructTypeNewFunc(file, service)
	for _, method := range service.Methods {
		generateServiceStructMethodHandler(file, method)
	}
	return file
}

func generateServiceStructType(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
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

func generateServiceStructTypeNewFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", "context")
	file.AddImport("", "github.com/wlMalk/goms/goms/errors")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	serviceName := strings.ToUpperFirst(service.Name)
	lowerServiceName := strings.ToLowerFirst(service.Name)
	file.Pf("func Endpoints(h interface{}, validatorsGetter interface{}, middlewareGetter interface{}) %s {", serviceName)
	file.Pf("handler := &endpointsHandler{")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: endpoint.Endpoint(func(ctx context.Context, req interface{}) (res interface{}, err error) {", lowerMethodName)
		file.Pf("return nil, errors.ErrMethodNotImplemented(\"%s\", \"%s\")", serviceName, methodName)
		file.Pf("}),")
	}
	file.Pf("}")
	file.Pf("s := &%s {", lowerServiceName)
	file.P("endpoints: handler,")
	file.P("handler:   handler,")
	file.Pf("}")
	file.Pf("")
	generateServiceMethodsRegisteration(file, service)
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodsRegisteration(file *files.GoFile, service *types.Service) {
	for _, method := range service.Methods {
		generateTypeSwitchForMethodHandler(file, method)
		if len(method.Arguments) > 0 && method.Options.Generate.Validator {
			generateMethodRequestValidatorMiddleware(file, method)
		}
		if method.Options.Generate.Middleware {
			generateMiddlewareCheckerForEndpoint(file, method)
		}
	}
	if service.Options.Generate.Middleware {
		generateMiddlewareCheckerForService(file, service)
	}
	generateEndpointsPacker(file, service)
}

func generateTypeSwitchForMethodHandler(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	file.AddImport("", method.Service.ImportPath, "/service/handlers/converters")
	file.Pf("switch t := h.(type) {")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("case handlers.%sHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.%sRequestResponseHandlerToEndpoint(converters.%sRequestHandlerTo%sRequestResponseHandler(converters.%sHandlerTo%sRequestHandler(handlers.%sHandlerFunc(t.%s))))", lowerMethodName, methodName, methodName, methodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.%sRequestResponseHandlerToEndpoint(converters.%sRequestHandlerTo%sRequestResponseHandler(handlers.%sRequestHandlerFunc(t.%s)))", lowerMethodName, methodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestResponseHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.%sRequestResponseHandlerToEndpoint(handlers.%sRequestResponseHandlerFunc(t.%s))", lowerMethodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sEndpointHandler:", methodName)
	file.Pf("s.endpoints.%s = t.%s", lowerMethodName, methodName)
	file.Pf("}")
	file.Pf("")
}

func generateMethodRequestValidatorMiddleware(file *files.GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "/service/requests")
	file.Pf("if t, ok := validatorsGetter.(interface {")
	file.Pf("Validate%s(ctx context.Context, req *requests.%sRequest) error", methodName, methodName)
	file.Pf("}); ok {")
	file.Pf("s.endpoints.%s = endpoint.Middleware(func(next endpoint.Endpoint) endpoint.Endpoint {", lowerMethodName)
	file.Pf("return endpoint.Endpoint(func(ctx context.Context, req interface{}) (interface{}, error) {")
	file.Pf("if err := t.Validate%s(ctx, req.(*requests.%sRequest)); err != nil {", methodName, methodName)
	file.Pf("return nil, err")
	file.Pf("}")
	file.Pf("return next(ctx, req)")
	file.Pf("})")
	file.Pf("})(s.endpoints.%s)", lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func generateMiddlewareCheckerForEndpoint(file *files.GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("if t, ok := middlewareGetter.(interface{ %sMiddleware(e endpoint.Endpoint) endpoint.Endpoint }); ok {", methodName)
	file.Pf("s.endpoints.%s = t.%sMiddleware(s.endpoints.%s)", lowerMethodName, methodName, lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func generateMiddlewareCheckerForService(file *files.GoFile, service *types.Service) {
	file.Pf("if t, ok := middlewareGetter.(interface{ Middleware(h handlers.EndpointHandler) handlers.EndpointHandler }); ok {")
	file.Pf("s.handler = t.Middleware(s.endpoints)")
	file.Pf("}")
	file.Pf("")
}

func generateOuterMiddlewareCheckerForEndpoint(file *files.GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("if t, ok := h.(interface {")
	file.Pf("Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint", methodName)
	file.Pf("}); ok {")
	file.Pf("%s = t.Outer%sMiddleware(%s)", lowerMethodName, methodName, lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func generateServiceStructMethodHandler(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("func (s *endpointsHandler) %s(ctx context.Context, req interface{}) (res interface{}, err error) {", methodName)
	file.Pf("return s.%s(ctx, req)", lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func generateEndpointsPacker(file *files.GoFile, service *types.Service) {
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
		generateOuterMiddlewareCheckerForEndpoint(file, method)
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
