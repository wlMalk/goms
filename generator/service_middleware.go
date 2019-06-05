package generator

import (
	"github.com/wlMalk/goms/parser/types"
)

func generateServiceMiddlewareFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateServiceMiddlewareType(file, service)
	return file
}

func generateServiceMiddlewareType(file *GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.P("type Middleware func(handlers.RequestResponseHandler) handlers.RequestResponseHandler")
	file.P("")
	file.P("func Chain(outer Middleware, others ...Middleware) Middleware {")
	file.P("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
	file.P("for i := len(others) - 1; i >= 0; i-- {")
	file.P("next = others[i](next)")
	file.P("}")
	file.P("return outer(next)")
	file.P("}")
	file.P("}")
}
