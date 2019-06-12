package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateRecoveringMiddlewareFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateRecoveringMiddlewareStruct(file, service)
	generateRecoveringMiddlewareNewFunc(file, service)
	for _, method := range service.Methods {
		generateRecoveringMiddlewareMethodFunc(file, method)
	}
	return file
}

func generateRecoveringMiddlewareStruct(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.P("type recoveringMiddleware struct {")
	file.P("next  handlers.RequestResponseHandler")
	file.P("}")
	file.P("")
}

func generateRecoveringMiddlewareNewFunc(file *files.GoFile, service *types.Service) {
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.P("func RecoveringMiddleware() RequestResponseMiddleware {")
	file.P("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
	file.P("return &recoveringMiddleware{")
	file.P("next: next,")
	file.P("}")
	file.P("}")
	file.P("}")
	file.P("")
}

func generateRecoveringMiddlewareMethodFunc(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", "fmt")
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
	file.Pf("func (m *recoveringMiddleware) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if method.Options.Generate.Middleware && method.Options.Generate.Recovering {
		file.Pf("defer func() {")
		file.Pf("if r:=recover(); r!=nil{")
		file.P("err = fmt.Errorf(\"%v\", r)")
		file.Pf("}")
		file.Pf("}()")
	}
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	file.Pf("return m.next.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}
