package generator

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateServiceStructType(file *GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("type %s struct {", serviceName)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%sHandler handlers.%sRequestResponseHandler", lowerMethodName, methodName)
	}
	file.Pf("}")
	file.Pf("")
}

func generateServiceStructTypeNewFunc(file *GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func New() *%s {", serviceName)
	file.Pf("return &%s{", serviceName)
	// Import errors
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		args := []string{"ctx context.Context"}
		if len(method.Arguments) > 0 {
			// Import requests
			args = append(args, "req *requests."+methodName+"Request")
		}
		results := []string{"err error"}
		if len(method.Results) > 0 {
			// Import responses
			results = append([]string{"res *responses." + methodName + "Response"}, results...)
		}
		file.Pf("%sHandler: handlers.%sRequestResponseHandlerFunc(func(%s) (%s) {", lowerMethodName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
		file.Pf("return nil, errors.Err%sMethodNotImplemented", methodName)
		file.Pf("}),")
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}

func generateServiceStructRegisterFunc(file *GoFile, service *types.Service) {
	file.Pf("func (s *%s) Register(h interface{}) {", strings.ToUpperFirst(service.Name))
	for _, method := range service.Methods {
		generateTypeSwitchForMethodHandler(file, method)
	}
	file.Pf("}")
	file.Pf("")
}

func generateTypeSwitchForMethodHandler(file *GoFile, method *types.Method) {
	file.Pf("switch t := h.(type) {")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	file.Pf("case func(%s) (%s):", strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("s.%sHandler = converters.%sRequestResponseHandlerTo%sRequestHandler(converters.%sRequestHandlerTo%sHandler(handlers.%sHandlerFunc(t)))", lowerMethodName, methodName, methodName, methodName, methodName, methodName)
	args = []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		// Import requests
		args = append(args, "req *requests."+methodName+"Request")
	}
	file.Pf("case func(%s) (%s):", strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("s.%sHandler = converters.%sRequestResponseHandlerTo%sRequestHandler(handlers.%sRequestHandlerFunc(t))", lowerMethodName, methodName, methodName, methodName)
	results = []string{"err error"}
	if len(method.Results) > 0 {
		// Import responses
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("case func(%s) (%s):", strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("s.%sHandler = handlers.%sRequestResponseHandlerFunc(t)", lowerMethodName, methodName)
	file.Pf("case handlers.%sHandlerFunc:", methodName)
	file.Pf("s.%sHandler = converters.%sRequestResponseHandlerTo%sRequestHandler(converters.%sRequestHandlerTo%sHandler(t))", lowerMethodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestHandlerFunc:", methodName)
	file.Pf("s.%sHandler = converters.%sRequestResponseHandlerTo%sRequestHandler(t)", lowerMethodName, methodName, methodName)
	file.Pf("case handlers.%sRequestResponseHandlerFunc:", methodName)
	file.Pf("s.%sHandler = t", lowerMethodName)
	file.Pf("case handlers.%sHandler:", methodName)
	file.Pf("s.%sHandler = converters.%sRequestResponseHandlerTo%sRequestHandler(converters.%sRequestHandlerTo%sHandler(handlers.%sHandlerFunc(t.%s)))", lowerMethodName, methodName, methodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestHandler:", methodName)
	file.Pf("s.%sHandler = converters.%sRequestResponseHandlerTo%sRequestHandler(handlers.%sRequestHandlerFunc(t.%s))", lowerMethodName, methodName, methodName, methodName, methodName)
	file.Pf("case handlers.%sRequestResponseHandler:", methodName)
	file.Pf("s.%sHandler = handlers.%sRequestResponseHandlerFunc(t.%s)", lowerMethodName, methodName, methodName)
	file.Pf("}")
	file.Pf("")
}

func generateServiceStructMethodHandler(file *GoFile, service string, method *types.Method) {
	serviceName := strings.ToUpperFirst(service)
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := []string{"ctx context.Context"}
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		// Import requests
		args = append(args, "req *requests."+methodName+"Request")
		argsInCall = append(argsInCall, "req")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		// Import responses
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func (s *%s) %s(%s) (%s) {", serviceName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return s.%sHandler.%s(%s)", lowerMethodName, methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}

func generateServiceStructFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateServiceStructType(file, service)
	generateServiceStructTypeNewFunc(file, service)
	generateServiceStructRegisterFunc(file, service)
	for _, method := range service.Methods {
		generateServiceStructMethodHandler(file, service.Name, method)
	}
	return file
}
