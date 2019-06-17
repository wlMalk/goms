package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ValidatingValidatorsTypes(file file.File, service types.Service) {
	file.P("type (")
	for _, method := range helpers.GetMethodsWithValidatingEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		file.Pf("%sValidator interface{", methodName)
		file.Pf("Validate%s(ctx context.Context, req *requests.%sRequest) (err error)", methodName, methodName)
		file.Pf("}")
	}
	file.P(")")
	file.P("")
}

func ValidatingMiddlewareStruct(file file.File, service types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.P("type validatingMiddleware struct {")
	for _, method := range helpers.GetMethodsWithValidatingEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%sValidator func(ctx context.Context, req *requests.%sRequest) (err error)", lowerMethodName, methodName)
	}
	file.P("")
	file.P("next handlers.RequestResponseHandler")
	file.P("}")
	file.P("")
}

func ValidatingMiddlewareNewFunc(file file.File, service types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.P("func ValidatingMiddleware(validator interface{}) RequestResponseMiddleware {")
	file.P("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
	file.Pf("m := &validatingMiddleware{next: next}")
	for _, method := range helpers.GetMethodsWithValidatingEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("if v, ok := validator.(%sValidator); ok {", methodName)
		file.Pf("m.%sValidator = v.Validate%s", lowerMethodName, methodName)
		file.Pf("}")
	}
	file.Pf("return m")
	file.P("}")
	file.P("}")
	file.P("")
}

func ValidatingMiddlewareMethodFunc(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func (m *validatingMiddleware) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if method.Options.Generate.Middleware && method.Options.Generate.Validating {
		file.Pf("if m.%sValidator != nil {", lowerMethodName)
		file.Pf("err = m.%sValidator(ctx, req)", lowerMethodName)
		file.Pf("if err != nil {")
		file.Pf("return")
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
