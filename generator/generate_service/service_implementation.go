package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateServiceImplementationFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, false, false)
	generateServiceImplementationStruct(file, service)
	generateServiceImplementationStructNewFunc(file, service)
	for _, method := range helpers.GetMethodsWithMethodStubsEnabled(service) {
		generateServiceMethodImplementation(file, service.Name, method)
	}
	return file
}

func GenerateServiceImplementationValidatorFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, false, false)
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("type %sValidator struct {}", serviceName)
	file.Pf("")
	file.Pf("func NewValidator() *%sValidator {", serviceName)
	file.Pf("return &%sValidator{}", serviceName)
	file.Pf("}")
	file.Pf("")
	for _, method := range helpers.GetMethodsWithValidatingEnabled(service) {
		generateServiceMethodImplementationValidateFunc(file, method)
	}
	return file
}

func GenerateServiceImplementationMiddlewareFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, false, false)
	for _, method := range helpers.GetMethodsWithMiddlewareEnabled(service) {
		generateServiceMethodImplementationMiddleware(file, method)
		generateServiceMethodImplementationOuterMiddleware(file, method)
	}
	if service.Options.Generate.Middleware {
		generateServiceImplementationMiddleware(file, service)
	}
	return file
}

func GenerateCachingKeyerFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, false, false)
	generateCachingMiddlewareCacheKeyerType(file, service)
	generateCachingMiddlewareKeyerNewFunc(file, service)
	for _, method := range helpers.GetMethodsWithCachingEnabled(service) {
		if len(method.Arguments) > 0 && len(method.Results) > 0 {
			generateCachingMiddlewareKeyerMethodFunc(file, method)
		}
	}
	return file
}

func generateServiceImplementationStruct(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("type %s struct {}", serviceName)
	file.Pf("")
}

func generateServiceImplementationStructNewFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func New() *%s {", serviceName)
	file.Pf("return &%s{}", serviceName)
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementation(file *files.GoFile, service string, method *types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	file.Pf("func (s *%s) %s(%s) (%s) {", serviceName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Cf("TODO: Implement %s method", methodName)
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementationMiddleware(file *files.GoFile, method *types.Method) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(method.Service.Name)
	file.Pf("func (s *%s) %sMiddleware(e endpoint.Endpoint) endpoint.Endpoint {", serviceName, methodName)
	file.Cf("TODO: Wrap %s middleware around e", methodName)
	file.Pf("return e")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementationOuterMiddleware(file *files.GoFile, method *types.Method) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(method.Service.Name)
	file.Pf("func (s *%s) Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint {", serviceName, methodName)
	file.Cf("TODO: Wrap %s middleware around e", methodName)
	file.Pf("return e")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMethodImplementationValidateFunc(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(method.Service.Name)
	file.Pf("func (v *%sValidator) Validate%s(ctx context.Context, req requests.%sRequest) error {", serviceName, methodName, methodName)
	file.Cf("TODO: Validate %s request", methodName)
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateServiceImplementationMiddleware(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func (s *%s) Middleware(h handlers.RequestResponseHandler) handlers.RequestResponseHandler {", serviceName)
	file.Cf("TODO: Wrap %s middleware around h", serviceName)
	file.Pf("return h")
	file.Pf("}")
	file.Pf("")
}

func generateCachingMiddlewareCacheKeyerType(file *files.GoFile, service *types.Service) {
	file.Pf("type CacheKeyer struct {")
	file.Pf("}")
	file.Pf("")
}

func generateCachingMiddlewareKeyerNewFunc(file *files.GoFile, service *types.Service) {
	file.P("func NewCacheKeyer() *CacheKeyer{")
	file.P("return &CacheKeyer{}")
	file.P("}")
	file.Pf("")
}

func generateCachingMiddlewareKeyerMethodFunc(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func (ck *CacheKeyer) %s(ctx context.Context, req *requests.%sRequest) (keys []interface{}, ok bool) {", methodName, methodName)
	file.Cf("TODO: Implement %s cache keyer method", methodName)
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}
