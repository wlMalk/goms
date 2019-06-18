package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func MethodHandlers(file file.File, service types.Service, method types.Method) error {
	MethodHandlerTypes(file, service, method)
	MethodHandlerFuncTypes(file, service, method)
	MethodHandlerFuncHandlers(file, service, method)
	return nil
}

func ServiceHandlerTypes(file file.File, service types.Service) error {
	file.Pf("type (")
	file.Pf("Handler interface {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sHandler", methodName)
	}
	file.Pf("}")
	file.Pf("RequestResponseHandler interface {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sRequestResponseHandler", methodName)
	}
	file.Pf("}")
	file.Pf("EndpointHandler interface {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sEndpointHandler", methodName)
	}
	file.Pf("}")
	file.Pf(")")
	file.Pf("")
	return nil
}

func MethodHandlerTypes(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	file.Pf("type (")
	file.Pf("%sHandler interface {", methodName)
	file.Pf("%s(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("}")
	args = []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results = []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("%sRequestResponseHandler interface {", methodName)
	file.Pf("%s(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("}")
	file.Pf("%sEndpointHandler interface {", methodName)
	file.Pf("%s(ctx context.Context, req interface{}) (res interface{}, err error)", methodName)
	file.Pf("}")
	file.Pf(")")
	file.Pf("")
	return nil
}

func MethodHandlerFuncTypes(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	file.Pf("type (")

	file.Pf("%sHandlerFunc func(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))

	args = []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results = []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("%sRequestResponseHandlerFunc func(%s) (%s)", methodName, strs.Join(args, ", "), strs.Join(results, ", "))

	file.Pf("%sEndpointHandlerFunc func(ctx context.Context, req interface{}) (res interface{}, err error)", methodName)

	file.Pf(")")
	file.Pf("")
	return nil
}

func MethodHandlerFuncHandlers(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	argsInCall := append([]string{"ctx"}, helpers.GetMethodArgumentsInCall(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	file.Pf("func (f %sHandlerFunc) %s(%s) (%s) {", methodName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return f(%s)", strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
	args = []string{"ctx context.Context"}
	argsInCall = []string{"ctx"}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
		argsInCall = append(argsInCall, "req")
	}
	results = []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func (f %sRequestResponseHandlerFunc) %s(%s) (%s) {", methodName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return f(%s)", strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
	file.Pf("func (f %sEndpointHandlerFunc) %s(ctx context.Context, req interface{}) (res interface{}, err error) {", methodName, methodName)
	file.Pf("return f(ctx, req)")
	file.Pf("}")
	file.Pf("")
	return nil
}
