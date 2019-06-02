package generator

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHandlersFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateServiceHandlerTypes(file, service)
	for _, method := range service.Methods {
		generateMethodHandlers(file, method)
	}
	return file
}

func generateMethodHandlers(file *GoFile, method *types.Method) {
	generateMethodHandlerTypes(file, method)
	generateMethodHandlerFuncTypes(file, method)
	generateMethodHandlerFuncHandlers(file, method)
}

func generateServiceHandlerTypes(file *GoFile, service *types.Service) {
	file.Pf("type (")
	file.Pf("Handler interface {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sHandler", methodName)
	}
	file.Pf("}")
	file.Pf("RequestHandler interface {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sRequestHandler", methodName)
	}
	file.Pf("}")
	file.Pf("RequestResponseHandler interface {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sRequestResponseHandler", methodName)
	}
	file.Pf("}")
	file.Pf(")")
	file.Pf("")
}

func generateMethodHandlerTypes(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	file.Pf("type (")
	file.Pf("%sHandler interface {", methodName)
	file.Pf("%s(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("}")
	args = []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	file.Pf("%sRequestHandler interface {", methodName)
	file.Pf("%s(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("}")
	results = []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("%sRequestResponseHandler interface {", methodName)
	file.Pf("%s(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("}")
	file.Pf(")")
	file.Pf("")
}

func generateMethodHandlerFuncTypes(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	file.Pf("type (")

	file.Pf("%sHandlerFunc func(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))

	args = []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	file.Pf("%sRequestHandlerFunc func(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))

	results = []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("%sRequestResponseHandlerFunc func(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))

	file.Pf(")")
	file.Pf("")
}

func generateMethodHandlerFuncHandlers(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	argsInCall := append([]string{"ctx"}, getMethodArgumentsInCall(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	file.Pf("func (f %sHandlerFunc) %s(%s) (%s) {", methodName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return f(%s)", strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
	args = []string{"ctx context.Context"}
	argsInCall = []string{"ctx"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
		argsInCall = append(argsInCall, "req")
	}
	file.Pf("func (f %sRequestHandlerFunc) %s(%s) (%s) {", methodName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return f(%s)", strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
	results = []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func (f %sRequestResponseHandlerFunc) %s(%s) (%s) {", methodName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return f(%s)", strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}
