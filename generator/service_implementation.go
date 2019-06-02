package generator

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateServiceImplementationFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, false, false)
	generateServiceImplementationStruct(file, service)
	generateServiceImplementationStructNewFunc(file, service)
	for _, method := range service.Methods {
		generateServiceMethodImplementation(file, service.Name, method)
		generateServiceMethodImplementationMiddleware(file, service.Name, method)
	}
	return file
}

func generateServiceImplementationStruct(file *GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("type %s struct {}", serviceName)
	file.Pf("")
}

func generateServiceImplementationStructNewFunc(file *GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func New() *%s {", serviceName)
	file.Pf("return &%s{}", serviceName)
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementation(file *GoFile, service string, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	file.Pf("func (s *%s) %s(%s) (%s) {", serviceName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Cf("TODO: Implement %s method", methodName)
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementationMiddleware(file *GoFile, service string, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service)
	file.Pf("func (s *%s) %sMiddleware() []interface{} {", serviceName, methodName)
	file.Cf("TODO: Add %sMiddleware", methodName)
	file.Pf("return []interface{}{}")
	file.Pf("}")
	file.Pf("")
}
