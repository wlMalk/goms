package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/constants"
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func LoggingMiddlewareStructs(file file.File, service types.Service) error {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	if helpers.HasLoggeds(service) {
		file.Pf("type loggingMiddleware struct {")
		file.Pf("next handlers.RequestResponseHandler")
		file.Pf("}")
		file.Pf("")
	}
	if helpers.HasLoggedErrors(service) {
		file.Pf("type errorLoggingMiddleware struct {")
		file.Pf("next handlers.RequestResponseHandler")
		file.Pf("}")
		file.Pf("")
	}
	return nil
}

func LoggingMiddlewareNewFunc(file file.File, service types.Service) error {
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	if helpers.HasLoggeds(service) {
		file.Pf("func LoggingMiddleware() RequestResponseMiddleware {")
		file.Pf("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
		file.Pf("return &loggingMiddleware{}")
		file.Pf("}")
		file.Pf("}")
		file.Pf("")
	}
	if helpers.HasLoggedErrors(service) {
		file.Pf("func ErrorLoggingMiddleware() RequestResponseMiddleware {")
		file.Pf("return func(next handlers.RequestResponseHandler) handlers.RequestResponseHandler {")
		file.Pf("return &errorLoggingMiddleware{}")
		file.Pf("}")
		file.Pf("}")
		file.Pf("")
	}
	return nil
}

func LoggingMiddlewareMethodHandler(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", "github.com/wlMalk/goms/goms/log")
	methodName := strings.ToUpperFirst(method.Name)
	serviceName := strings.ToUpperFirst(service.Name)
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
	file.Pf("func (m *loggingMiddleware) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if method.Generate.Has(constants.MethodGenerateMiddlewareFlag, constants.MethodGenerateLoggingFlag) && (helpers.HasLoggedArguments(method) || helpers.HasLoggedResults(method)) {
		file.Pf("defer func() {")
		file.Pf("if err == nil {")
		file.Pf("log.Info(ctx,")
		file.Pf("\"service\", \"%s\",", helpers.GetName(serviceName, service.Alias))
		file.Pf("\"method\", \"%s\",", helpers.GetName(methodName, method.Alias))
		if helpers.HasLoggedArguments(method) {
			file.Pf("\"request\", log%sRequest{", methodName)
			for _, arg := range helpers.GetLoggedArgumentsForMethod(method) {
				argName := strings.ToUpperFirst(arg.Name)
				file.Pf("%s: req.%s,", argName, argName)
			}
			for _, arg := range helpers.GetLoggedArgumentsLenForMethod(method) {
				argName := strings.ToUpperFirst(arg.Name)
				file.Pf("Len%s: len(req.%s),", argName, argName)
			}
			file.Pf("},")
		}
		if helpers.HasLoggedResults(method) {
			file.Pf("\"response\", log%sResponse{", methodName)
			for _, field := range helpers.GetLoggedResultsForMethod(method) {
				fieldName := strings.ToUpperFirst(field.Name)
				file.Pf("%s: res.%s,", fieldName, fieldName)
			}
			for _, field := range helpers.GetLoggedResultsLenForMethod(method) {
				fieldName := strings.ToUpperFirst(field.Name)
				file.Pf("Len%s: len(res.%s),", fieldName, fieldName)
			}
			file.Pf("},")
		}
		file.Pf(")")
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
	return nil
}

func ErrorLoggingMiddlewareMethodHandler(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", "github.com/wlMalk/goms/goms/log")
	methodName := strings.ToUpperFirst(method.Name)
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
	file.Pf("func (m *errorLoggingMiddleware) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if method.Generate.Has(constants.MethodGenerateMiddlewareFlag, constants.MethodGenerateLoggingFlag) && !method.Options.Logging.IgnoreError {
		file.Pf("defer func() {")
		file.Pf("if err != nil {")
		file.Pf("log.Error(ctx, \"error\", err)")
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
	return nil
}

func LoggingMiddlewareTypes(file file.File, service types.Service) error {
	helpers.AddTypesImports(file, service)
	methods := helpers.GetMethodsWithLoggingEnabled(service)
	if len(methods) > 0 {
		file.Pf("type (")
		for _, method := range methods {
			methodName := strings.ToUpperFirst(method.Name)
			if method.Generate.Has(constants.MethodGenerateMiddlewareFlag, constants.MethodGenerateLoggingFlag) {
				if helpers.HasLoggedArguments(method) {
					file.Pf("log%sRequest struct {", methodName)
					for _, arg := range helpers.GetLoggedArgumentsForMethod(method) {
						argName := strings.ToUpperFirst(arg.Name)
						argSpecialName := helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias)
						file.Pf("%s %s `json:\"%s\"`", argName, arg.Type.GoType(), argSpecialName)
					}
					for _, arg := range helpers.GetLoggedArgumentsLenForMethod(method) {
						argName := strings.ToUpperFirst(arg.Name)
						argSpecialName := helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias)
						file.Pf("Len%s %s `json:\"len(%s)\"`", argName, arg.Type.GoType(), argSpecialName)
					}
					file.Pf("}")
				}
				if helpers.HasLoggedResults(method) {
					file.Pf("log%sResponse struct {", methodName)
					for _, field := range helpers.GetLoggedResultsForMethod(method) {
						fieldName := strings.ToUpperFirst(field.Name)
						fieldSpecialName := helpers.GetName(strings.ToLowerFirst(field.Name), field.Alias)
						file.Pf("%s %s `json:\"%s\"`", fieldName, field.Type.GoType(), fieldSpecialName)
					}
					for _, field := range helpers.GetLoggedResultsLenForMethod(method) {
						fieldName := strings.ToUpperFirst(field.Name)
						fieldSpecialName := helpers.GetName(strings.ToLowerFirst(field.Name), field.Alias)
						file.Pf("Len%s %s `json:\"len(%s)\"`", fieldName, field.Type.GoType(), fieldSpecialName)
					}
					file.Pf("}")
				}
			}
		}
		file.Pf(")")
		file.Pf("")
	}
	return nil
}
