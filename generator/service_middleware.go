package generator

import (
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateServiceMiddlewareType(file *GoFile, service *types.Service) {
	file.Pf("type Middleware func(handlers.RequestResponseHandler) handlers.RequestResponseHandler")
	file.Pf("")
}

func generateMethodMiddlewareType(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("type %sMiddleware func(handlers.%sRequestResponseHandler) handlers.%sRequestResponseHandler", methodName, methodName, methodName)
	file.Pf("")
}

func generateServiceMiddlewareFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateServiceMiddlewareType(file, service)
	for _, method := range service.Methods {
		generateMethodMiddlewareType(file, method)
	}
	return file
}
