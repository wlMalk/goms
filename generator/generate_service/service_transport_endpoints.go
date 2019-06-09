package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
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
	file.Pf("handler handlers.RequestResponseHandler")
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
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.Pf("func Endpoints(h interface{}) *%s {", serviceName)
	file.Pf("handler := &endpointsHandler{")
	for _, method := range service.Methods {
		lowerMethodName := strings.ToLowerFirst(method.Name)
		methodNameSnake := strings.ToSnakeCase(method.Name)
		file.Pf("%s: endpoint.Endpoint(func(ctx context.Context, req interface{}) (res interface{}, err error) {", lowerMethodName)
		file.Pf("return nil, errors.ErrMethodNotImplemented(\"%s\", \"%s\")", serviceNameSnake, methodNameSnake)
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
		generateMethodRequestValidatorMiddleware(file, method)
		generateMiddlewareCheckerForEndpoint(file, method)
	}
	generateMiddlewareCheckerForService(file, service)

	generateEndpointsPacker(file, service)
	file.Pf("")
}

func generateTypeSwitchForMethodHandler(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	file.AddImport("", method.Service.ImportPath, "/service/handlers/converters")
	file.Pf("switch t := h.(type) {")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("case handlers.%sHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.EndpointTo%sRequestResponseHandler(converters.%sRequestResponseHandlerTo%sRequestHandler(converters.%sRequestHandlerTo%sHandler(handlers.%sHandlerFunc(t.%s))))", lowerMethodName, methodName, methodName, methodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.EndpointTo%sRequestResponseHandler(converters.%sRequestResponseHandlerTo%sRequestHandler(handlers.%sRequestHandlerFunc(t.%s)))", lowerMethodName, methodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestResponseHandler:", methodName)
	file.Pf("s.endpoints.%s = converters.EndpointTo%sRequestResponseHandler(handlers.%sRequestResponseHandlerFunc(t.%s))", lowerMethodName, methodName, methodName, methodName)
	file.Pf("}")
	file.Pf("")
}

func generateMethodRequestValidatorMiddleware(file *files.GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("if t, ok := h.(interface {")
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
	file.Pf("if t, ok := h.(interface{ %sMiddleware(e endpoint.Endpoint) endpoint.Endpoint }); ok {", methodName)
	file.Pf("s.endpoints.%s = t.%sMiddleware(s.endpoints.%s)", lowerMethodName, methodName, lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func generateMiddlewareCheckerForService(file *files.GoFile, service *types.Service) {
	file.Pf("if t, ok := h.(interface{ Middleware(h handlers.RequestResponseHandler) handlers.RequestResponseHandler }); ok {")
	file.Pf("s.handler = t.Middleware(s.endpoints)")
	file.Pf("}")
	file.Pf("")
}

func generateOuterMiddlewareCheckerForEndpoint(file *files.GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.Pf("if t, ok := h.(interface{ Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint }); ok {", methodName)
	file.Pf("s.endpoints.%s = t.Outer%sMiddleware(s.endpoints.%s)", lowerMethodName, methodName, lowerMethodName)
	file.Pf("}")
	file.Pf("")
}

func generateServiceStructMethodHandler(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func (s *endpointsHandler) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	req := "nil"
	if len(method.Arguments) > 0 {
		req = "req"
	}
	if len(method.Results) > 0 {
		file.Pf("resp, err := s.%s(ctx, %s)", lowerMethodName, req)
	} else {
		file.Pf("_, err = s.%s(ctx, %s)", lowerMethodName, req)
	}
	if len(method.Results) > 0 {

		file.Pf("if err != nil {")
		file.Pf("return nil, err")
		file.Pf("}")
		file.Pf("if resp != nil {")
		file.Pf("return resp.(*responses.%sResponse), nil", methodName)
		file.Pf("}")
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
	file.Pf("")
}

func generateEndpointsPacker(file *files.GoFile, service *types.Service) {
	file.Pf("var (")
	serviceName := strings.ToUpperFirst(service.Name)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s endpoint.Endpoint = converters.EndpointTo%sRequestResponseHandler(handlers.%sRequestResponseHandlerFunc(s.handler.%s))", lowerMethodName, methodName, methodName, methodName)
	}
	file.Pf(")")
	file.Pf("")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("if t, ok := h.(interface {")
		file.Pf("Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint", methodName)
		file.Pf("}); ok {")
		file.Pf("%s = t.Outer%sMiddleware(%s)", lowerMethodName, methodName, lowerMethodName)
		file.Pf("}")
		file.Pf("")
	}
	file.Pf("endpoints := &%s{", serviceName)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: %s,", methodName, lowerMethodName)
	}
	file.Pf("}")
	file.Pf("")
	file.Pf("return endpoints")
	file.Pf("")
}
