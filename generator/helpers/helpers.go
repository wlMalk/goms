package helpers

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GetMethodArguments(args []*types.Argument) []string {
	var a []string
	for _, arg := range args {
		a = append(a, fmt.Sprintf("%s %s", strings.ToLowerFirst(arg.Name), arg.Type.GoArgumentType()))
	}
	return a
}

func GetMethodResults(results []*types.Field) []string {
	var r []string
	for _, result := range results {
		r = append(r, fmt.Sprintf("%s %s", strings.ToLowerFirst(result.Name), result.Type.GoArgumentType()))
	}
	return r
}

func AddMethodImports(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
}

func GenerateFunc(file *files.GoFile, receiver string, name string, args []string, results []string, f func()) {
	file.P(GetFuncSignature(receiver, name, args, results))
	f()
	file.P("}")
	file.P("")
}

func GetFuncSignature(receiver string, name string, args []string, results []string) string {
	signature := "func "
	if receiver != "" {
		signature += "(" + receiver + ") "
	}
	signature += name + "(" + strs.Join(args, ", ") + ") "
	if len(results) > 0 {
		signature += "(" + strs.Join(results, ", ") + ") "
	}
	signature += "{"
	return signature
}

func GetExportedMethodSignature(receiver string, method *types.Method) string {
	return GetFuncSignature(receiver, strings.ToUpperFirst(method.Name), GetMethodArguments(method.Arguments), GetMethodResults(method.Results))
}

func GetUnexportedMethodSignature(receiver string, method *types.Method) string {
	return GetFuncSignature(receiver, strings.ToLowerFirst(method.Name), GetMethodArguments(method.Arguments), GetMethodResults(method.Results))
}

func GetFieldTagsString(tags map[string]string) string {
	var a []string
	for k, v := range tags {
		a = append(a, fmt.Sprintf("%s:\"%s\"", k, v))
	}
	return strs.Join(a, ",")
}

func GenerateStruct(file *files.GoFile, name string, fields []*types.Field) {
	if len(fields) == 0 {
		return
	}
	file.Pf("type %s struct {", name)
	for _, f := range fields {
		jsonName := GetName(f.Name, f.Alias)
		if len(f.Tags) == 0 {
			f.Tags = map[string]string{"json": strings.ToLowerFirst(jsonName)}
		} else {
			f.Tags["json"] = strings.ToLowerFirst(jsonName)
		}
		file.Pf("%s %s `%s`", f.Name, f.Type.GoType(), GetFieldTagsString(f.Tags))
	}
	file.P("}")
	file.P("")
}

func GenerateExportedStruct(file *files.GoFile, name string, fields []*types.Field) {
	GenerateStruct(file, strings.ToUpperFirst(name), fields)
}

func GenerateUnexportedStruct(file *files.GoFile, name string, fields []*types.Field) {
	GenerateStruct(file, strings.ToLowerFirst(name), fields)
}

func GetMethodArgumentsInCall(args []*types.Argument) (a []string) {
	for _, arg := range args {
		if arg.Type.IsVariadic {
			a = append(a, strings.ToLowerFirst(arg.Name)+"...")
		} else {
			a = append(a, strings.ToLowerFirst(arg.Name))
		}
	}
	return
}

func GetMethodArgumentsFromRequestInCall(args []*types.Argument) (a []string) {
	a = GetMethodArgumentsInCall(args)
	for i := range a {
		a[i] = "req." + strings.ToUpperFirst(a[i])
	}
	return
}

func GetResultsVars(results []*types.Field) (a []string) {
	for _, result := range results {
		a = append(a, strings.ToLowerFirst(result.Name))
	}
	return
}

func GetResultsVarsFromResponse(results []*types.Field) (a []string) {
	a = GetResultsVars(results)
	for i := range a {
		a[i] = "res." + strings.ToUpperFirst(a[i])
	}
	return
}

func GetName(name string, alias string) string {
	n := name
	if alias != "" {
		n = alias
	}
	return strings.ToLowerFirst(n)
}

func IsCachingEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Caching {
			return true
		}
	}
	return false
}

func IsLoggingEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Logging {
			return true
		}
	}
	return false
}

func IsMethodStubsEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.MethodStubs {
			return true
		}
	}
	return false
}

func IsValidatorsEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Validator {
			return true
		}
	}
	return false
}

func IsMiddlewareEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Middleware {
			return true
		}
	}
	return service.Options.Generate.Middleware
}

func IsHTTPServerEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Transports.HTTP.Server {
			return true
		}
	}
	return false
}

func IsHTTPClientEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Transports.HTTP.Client {
			return true
		}
	}
	return false
}

func IsGRPCServerEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Transports.GRPC.Server {
			return true
		}
	}
	return false
}

func IsGRPCClientEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Transports.GRPC.Client {
			return true
		}
	}
	return false
}

func IsTracingEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Tracing {
			return true
		}
	}
	return false
}

func IsMetricsEnabled(service *types.Service) bool {
	for _, method := range service.Methods {
		if method.Options.Metrics {
			return true
		}
	}
	return false
}

func FilteredMethods(methods []*types.Method, filter func(method *types.Method) bool) (ms []*types.Method) {
	for _, method := range methods {
		if filter(method) {
			ms = append(ms, method)
		}
	}
	return
}

func FilteredArguments(args []*types.Argument, filter func(arg *types.Argument) bool) (as []*types.Argument) {
	for _, arg := range args {
		if filter(arg) {
			as = append(as, arg)
		}
	}
	return
}

func FilteredFields(fields []*types.Field, filter func(field *types.Field) bool) (fs []*types.Field) {
	for _, field := range fields {
		if filter(field) {
			fs = append(fs, field)
		}
	}
	return
}

func GetMethodsWithCachingEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Caching
	})
}

func GetMethodsWithLoggingEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Logging
	})
}

func GetMethodsWithMethodStubsEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.MethodStubs
	})
}

func GetMethodsWithValidatorEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Validator
	})
}

func GetMethodsWithMiddlewareEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Middleware
	})
}

func GetMethodsWithHTTPServerEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Transports.HTTP.Server
	})
}

func GetMethodsWithHTTPClientEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Transports.HTTP.Server
	})
}

func GetMethodsWithHTTPEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Transports.HTTP.Server || method.Options.Transports.HTTP.Client
	})
}

func GetMethodsWithGRPCServerEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Transports.GRPC.Server
	})
}

func GetMethodsWithGRPCClientEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Transports.GRPC.Server
	})
}

func GetMethodsWithGRPCEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Transports.GRPC.Server || method.Options.Transports.GRPC.Client
	})
}

func GetMethodsWithTracingEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Tracing
	})
}

func GetMethodsWithMetricsEnabled(service *types.Service) (ms []*types.Method) {
	return FilteredMethods(service.Methods, func(method *types.Method) bool {
		return method.Options.Metrics
	})
}

func GetLoggedArgumentsForMethod(method *types.Method) (args []*types.Argument) {
	return FilteredArguments(method.Arguments, func(arg *types.Argument) bool {
		return !containsNamesAliases(method.Options.LoggingOptions.IgnoredArguments, arg.Name, arg.Alias)
	})
}

func GetLoggedResultsForMethod(method *types.Method) (results []*types.Field) {
	return FilteredFields(method.Results, func(field *types.Field) bool {
		return !containsNamesAliases(method.Options.LoggingOptions.IgnoredResults, field.Name, field.Alias)
	})
}

func GetLoggedArgumentsLenForMethod(method *types.Method) (args []*types.Argument) {
	return FilteredArguments(method.Arguments, func(arg *types.Argument) bool {
		return (arg.Type.IsMap || arg.Type.IsVariadic || arg.Type.IsSlice || arg.Type.IsBytes) && containsNamesAliases(method.Options.LoggingOptions.LenArguments, arg.Name, arg.Alias)
	})
}

func GetLoggedResultsLenForMethod(method *types.Method) (fields []*types.Field) {
	return FilteredFields(method.Results, func(field *types.Field) bool {
		return (field.Type.IsMap || field.Type.IsVariadic || field.Type.IsSlice || field.Type.IsBytes) && containsNamesAliases(method.Options.LoggingOptions.LenResults, field.Name, field.Alias)
	})
}

func containsNamesAliases(ss []string, name string, alias string) bool {
	for i := range ss {
		if strings.ToLower(ss[i]) == strings.ToLower(name) || (len(strs.TrimSpace(alias)) > 0 && strings.ToLower(ss[i]) == strings.ToLower(alias)) {
			return true
		}
	}
	return false
}
