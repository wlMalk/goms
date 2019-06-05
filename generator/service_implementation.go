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
	}
	return file
}

func generateServiceImplementationValidatorsFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, false, false)
	for _, method := range service.Methods {
		if len(method.Arguments) > 0 {
			generateServiceMethodImplementationValidateFunc(file, method)
		}
	}
	return file
}

func generateServiceImplementationMiddlewareFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, false, false)
	for _, method := range service.Methods {
		generateServiceMethodImplementationMiddleware(file, method)
		generateServiceMethodImplementationOuterMiddleware(file, method)
	}
	generateServiceImplementationMiddleware(file, service)
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

func generateServiceMethodImplementationMiddleware(file *GoFile, method *types.Method) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(method.Service.Name)
	file.Pf("func (s *%s) %sMiddleware(e endpoint.Endpoint) endpoint.Endpoint {", serviceName, methodName)
	file.Cf("TODO: Wrap %s middleware around e", methodName)
	file.Pf("return e")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementationOuterMiddleware(file *GoFile, method *types.Method) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(method.Service.Name)
	file.Pf("func (s *%s) Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint {", serviceName, methodName)
	file.Cf("TODO: Wrap %s middleware around e", methodName)
	file.Pf("return e")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementationValidateFunc(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(method.Service.Name)
	file.Pf("func (s *%s) Validate%s(ctx context.Context, req requests.%sRequest) error {", serviceName, methodName, methodName)
	file.Cf("TODO: Validate %s request", methodName)
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateServiceImplementationMiddleware(file *GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func (s *%s) Middleware(h handlers.RequestResponseHandler) handlers.RequestResponseHandler {", serviceName)
	file.Cf("TODO: Wrap %s middleware around h", serviceName)
	file.Pf("return h")
	file.Pf("}")
	file.Pf("")
}
