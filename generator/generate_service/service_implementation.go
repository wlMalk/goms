package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceImplementationStruct(file file.File, service types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("type %s struct {}", serviceName)
	file.Pf("")
}

func ServiceImplementationStructNewFunc(file file.File, service types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func New() *%s {", serviceName)
	file.Pf("return &%s{}", serviceName)
	file.Pf("}")
	file.Pf("")
}

func ServiceMethodImplementation(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	file.Pf("func (s *%s) %s(%s) (%s) {", serviceName, methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Cf("TODO: Implement %s method", methodName)
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}

func ServiceMethodImplementationMiddleware(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func (s *%s) %sMiddleware(e endpoint.Endpoint) endpoint.Endpoint {", serviceName, methodName)
	file.Cf("TODO: Wrap %s middleware around e", methodName)
	file.Pf("return e")
	file.Pf("}")
	file.Pf("")
}

func ServiceMethodImplementationOuterMiddleware(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func (s *%s) Outer%sMiddleware(e endpoint.Endpoint) endpoint.Endpoint {", serviceName, methodName)
	file.Cf("TODO: Wrap %s middleware around e", methodName)
	file.Pf("return e")
	file.Pf("}")
	file.Pf("")
}

func ServiceMethodImplementationValidateFunc(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func (v *%sValidator) Validate%s(ctx context.Context, req requests.%sRequest) error {", serviceName, methodName, methodName)
	file.Cf("TODO: Validate %s request", methodName)
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func ServiceImplementationMiddleware(file file.File, service types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func (s *%s) Middleware(h handlers.RequestResponseHandler) handlers.RequestResponseHandler {", serviceName)
	file.Cf("TODO: Wrap %s middleware around h", serviceName)
	file.Pf("return h")
	file.Pf("}")
	file.Pf("")
}

func CachingMiddlewareCacheKeyerType(file file.File, service types.Service) {
	file.Pf("type CacheKeyer struct {")
	file.Pf("}")
	file.Pf("")
}

func CachingMiddlewareKeyerNewFunc(file file.File, service types.Service) {
	file.P("func NewCacheKeyer() *CacheKeyer{")
	file.P("return &CacheKeyer{}")
	file.P("}")
	file.Pf("")
}

func CachingMiddlewareKeyerMethodFunc(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func (ck *CacheKeyer) %s(ctx context.Context, req *requests.%sRequest) (keys []interface{}, ok bool) {", methodName, methodName)
	file.Cf("TODO: Implement %s cache keyer method", methodName)
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}
