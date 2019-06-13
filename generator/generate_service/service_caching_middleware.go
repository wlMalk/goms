package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateCachingMiddlewareFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateCachingMiddlewareStruct(file, service)
	generateCachingMiddlewareCacheKeyerInterface(file, service)
	generateCachingMiddlewareNewFunc(file, service)
	for _, method := range service.Methods {
		generateCachingMiddlewareMethodFunc(file, method)
	}
	return file
}

func generateCachingMiddlewareStruct(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.AddImport("", "github.com/wlMalk/goms/goms/cache")
	file.AddImport("", "hash")
	file.P("type cachingMiddleware struct {")
	file.P("cache cache.Cache")
	file.P("keyer  cacheKeyer")
	file.P("hasher func() hash.Hash")
	file.P("next  handlers.RequestResponseHandler")
	file.P("}")
	file.P("")
}

func generateCachingMiddlewareCacheKeyerInterface(file *files.GoFile, service *types.Service) {
	file.Pf("type cacheKeyer interface {")
	for _, method := range service.Methods {
		if method.Options.Generate.Middleware && len(method.Arguments) > 0 && len(method.Results) > 0 && method.Options.Generate.Caching {
			methodName := strings.ToUpperFirst(method.Name)
			file.Pf("%s(ctx context.Context, req *requests.%sRequest) (keys []interface{}, ok bool)", methodName, methodName)
		}
	}
	file.Pf("}")
}

func generateCachingMiddlewareNewFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.AddImport("", "github.com/wlMalk/goms/goms/cache")
	file.AddImport("", "hash")
	file.P("func CachingMiddleware(cache cache.Cache, keyer cacheKeyer, hasher func() hash.Hash) RequestResponseMiddleware {")
	file.P("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
	file.P("return &cachingMiddleware{")
	file.P("cache:  cache,")
	file.P("keyer:  keyer,")
	file.P("hasher: hasher,")
	file.P("next:   next,")
	file.P("}")
	file.P("}")
	file.P("}")
	file.P("")
}

func generateCachingMiddlewareMethodFunc(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "github.com/wlMalk/goms/goms/cache")
	file.AddImport("", "github.com/wlMalk/goms/goms/log")
	methodName := strings.ToUpperFirst(method.Name)
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
	file.Pf("func (m *cachingMiddleware) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if method.Options.Generate.Middleware && len(method.Arguments) > 0 && len(method.Results) > 0 && method.Options.Generate.Caching {
		file.Pf("keys, ok := m.keyer.%s(ctx, req)", methodName)
		file.Pf("if ok {")
		file.Pf("key, err := cache.Key(m.hasher, keys...)")
		file.Pf("if err != nil {")
		file.Pf("log.Error(ctx, \"message\", err)")
		file.Pf("} else {")
		file.Pf("value, err := m.cache.Get(key)")
		file.Pf("if err == nil && value != nil {")
		file.Pf("res, ok := value.(*responses.%sResponse)", methodName)
		file.Pf("if ok {")
		file.Pf("return res, nil")
		file.Pf("}")
		file.Pf("}")
		file.Pf("defer func() {")
		file.Pf("m.cache.Set(key, res)")
		file.Pf("}()")
		file.Pf("}")
		file.Pf("}")
	}
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	file.Pf("return m.next.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}
