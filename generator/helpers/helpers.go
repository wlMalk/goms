package helpers

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/constants"
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func AddTypesImports(file file.File, service types.Service) {
	if len(service.Entities) > 0 || len(service.ArgumentsGroups) > 0 || len(service.Enums) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/types")
	}
}

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

func AddMethodImports(file *files.GoFile, method types.Method) {
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

func GetExportedMethodSignature(receiver string, method types.Method) string {
	return GetFuncSignature(receiver, strings.ToUpperFirst(method.Name), GetMethodArguments(method.Arguments), GetMethodResults(method.Results))
}

func GetUnexportedMethodSignature(receiver string, method types.Method) string {
	return GetFuncSignature(receiver, strings.ToLowerFirst(method.Name), GetMethodArguments(method.Arguments), GetMethodResults(method.Results))
}

func GetFieldTagsString(tags map[string][]string) string {
	var a []string
	for k, v := range tags {
		a = append(a, fmt.Sprintf("%s:\"%s\"", k, strs.Join(v, ",")))
	}
	return strs.Join(a, ",")
}

func GenerateStruct(file *files.GoFile, name string, fields []*types.Field) {
	if len(fields) == 0 {
		return
	}
	file.Pf("type %s struct {", name)
	for _, f := range fields {
		jsonName := GetName(strings.ToLowerFirst(f.Name), f.Alias)
		if len(f.Tags) == 0 {
			f.Tags = map[string][]string{"json": []string{strings.ToLowerFirst(jsonName)}}
		} else {
			f.Tags["json"] = []string{strings.ToLowerFirst(jsonName)}
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
	return n
}

func IsCachingEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateCachingFlag) {
			return true
		}
	}
	return false
}

func IsLoggingEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateLoggingFlag) {
			return true
		}
	}
	return false
}

func IsServerEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.HasAny(constants.MethodGenerateHTTPServerFlag, constants.MethodGenerateGRPCServerFlag) {
			return true
		}
	}
	return false
}

func IsRateLimitingEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateRateLimitingFlag) {
			return true
		}
	}
	return false
}

func IsCircuitBreakingEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateCircuitBreakingFlag) {
			return true
		}
	}
	return false
}

func HasLoggeds(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateLoggingFlag) && (HasLoggedArguments(method) || HasLoggedResults(method)) {
			return true
		}
	}
	return false
}

func HasLoggedErrors(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateLoggingFlag) && !method.Options.Logging.IgnoreError {
			return true
		}
	}
	return false
}

func IsMethodStubsEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateMethodStubsFlag) {
			return true
		}
	}
	return false
}

func IsValidatorsEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateValidatorsFlag) {
			return true
		}
	}
	return false
}

func IsValidatingEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateValidatingFlag) {
			return true
		}
	}
	return false
}

func IsMiddlewareEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateMiddlewareFlag) {
			return true
		}
	}
	return false
}

func IsRecoveringEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateRecoveringFlag) {
			return true
		}
	}
	return false
}

func IsHTTPEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.HasAny(constants.MethodGenerateHTTPServerFlag, constants.MethodGenerateHTTPClientFlag) {
			return true
		}
	}
	return false
}

func IsHTTPServerEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateHTTPServerFlag) {
			return true
		}
	}
	return false
}

func IsHTTPClientEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateHTTPClientFlag) {
			return true
		}
	}
	return false
}

func IsGRPCEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.HasAny(constants.MethodGenerateGRPCServerFlag, constants.MethodGenerateGRPCClientFlag) {
			return true
		}
	}
	return false
}

func IsGRPCServerEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateGRPCServerFlag) {
			return true
		}
	}
	return false
}

func IsGRPCClientEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateGRPCClientFlag) {
			return true
		}
	}
	return false
}

func IsTracingEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateTracingFlag) {
			return true
		}
	}
	return false
}

func IsMetricsEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.HasAny(
			constants.MethodGenerateFrequencyMetricFlag,
			constants.MethodGenerateLatencyMetricFlag,
			constants.MethodGenerateCounterMetricFlag,
		) {
			return true
		}
	}
	return false
}

func IsFrequencyMetricEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateFrequencyMetricFlag) {
			return true
		}
	}
	return false
}

func IsLatencyMetricEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateLatencyMetricFlag) {
			return true
		}
	}
	return false
}

func IsCounterMetricEnabled(service types.Service) bool {
	for _, method := range service.Methods {
		if method.Generate.Has(constants.MethodGenerateCounterMetricFlag) {
			return true
		}
	}
	return false
}

func FilteredMethods(methods []types.Method, filter func(method types.Method) bool) (ms []types.Method) {
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

func GetMethodsWithCachingEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateCachingFlag)
	})
}

func GetMethodsWithLoggingEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateLoggingFlag)
	})
}

func GetMethodsWithMethodStubsEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateMethodStubsFlag)
	})
}

func GetMethodsWithValidatorsEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateValidatorsFlag)
	})
}

func GetMethodsWithValidatingEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateValidatingFlag)
	})
}

func GetMethodsWithMiddlewareEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateMiddlewareFlag)
	})
}

func GetMethodsWithHTTPServerEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateHTTPServerFlag)
	})
}

func GetMethodsWithHTTPClientEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateHTTPServerFlag)
	})
}

func GetMethodsWithHTTPEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.HasAny(constants.MethodGenerateHTTPServerFlag, constants.MethodGenerateHTTPClientFlag)
	})
}

func GetMethodsWithGRPCServerEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateGRPCServerFlag)
	})
}

func GetMethodsWithGRPCClientEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateGRPCServerFlag)
	})
}

func GetMethodsWithGRPCEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.HasAny(constants.MethodGenerateGRPCServerFlag, constants.MethodGenerateGRPCClientFlag)
	})
}

func GetMethodsWithTracingEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.Has(constants.MethodGenerateTracingFlag)
	})
}

func GetMethodsWithMetricsEnabled(service types.Service) (ms []types.Method) {
	return FilteredMethods(service.Methods, func(method types.Method) bool {
		return method.Generate.HasAny(constants.MethodGenerateFrequencyMetricFlag, constants.MethodGenerateLatencyMetricFlag, constants.MethodGenerateCounterMetricFlag)
	})
}

func GetLoggedArgumentsForMethod(method types.Method) (args []*types.Argument) {
	return FilteredArguments(method.Arguments, func(arg *types.Argument) bool {
		return !containsNamesAliases(method.Options.Logging.IgnoredArguments, arg.Name, arg.Alias)
	})
}

func GetLoggedResultsForMethod(method types.Method) (results []*types.Field) {
	return FilteredFields(method.Results, func(field *types.Field) bool {
		return !containsNamesAliases(method.Options.Logging.IgnoredResults, field.Name, field.Alias)
	})
}

func GetLoggedArgumentsLenForMethod(method types.Method) (args []*types.Argument) {
	return FilteredArguments(method.Arguments, func(arg *types.Argument) bool {
		return (arg.Type.IsMap || arg.Type.IsVariadic || arg.Type.IsSlice || arg.Type.IsBytes) && containsNamesAliases(method.Options.Logging.LenArguments, arg.Name, arg.Alias)
	})
}

func GetLoggedResultsLenForMethod(method types.Method) (fields []*types.Field) {
	return FilteredFields(method.Results, func(field *types.Field) bool {
		return (field.Type.IsMap || field.Type.IsVariadic || field.Type.IsSlice || field.Type.IsBytes) && containsNamesAliases(method.Options.Logging.LenResults, field.Name, field.Alias)
	})
}

func HasLoggedArguments(method types.Method) bool {
	return (len(method.Arguments) > 0 && len(method.Options.Logging.IgnoredArguments) < len(method.Arguments)) ||
		len(method.Options.Logging.LenArguments) > 0
}

func HasLoggedResults(method types.Method) bool {
	return (len(method.Results) > 0 && len(method.Options.Logging.IgnoredResults) < len(method.Results)) ||
		len(method.Options.Logging.LenResults) > 0
}

func IsCachaeble(service types.Service) bool {
	for _, method := range service.Methods {
		if len(method.Arguments) > 0 && len(method.Results) > 0 && method.Generate.Has(constants.MethodGenerateCachingFlag) {
			return true
		}
	}
	return false
}

func IsValidatable(service types.Service) bool {
	for _, method := range service.Methods {
		if len(method.Arguments) > 0 && method.Generate.Has(constants.MethodGenerateValidatingFlag) {
			return true
		}
	}
	return false
}

func containsNamesAliases(ss []string, name string, alias string) bool {
	for i := range ss {
		if strings.ToLower(ss[i]) == strings.ToLower(name) || (len(strs.TrimSpace(alias)) > 0 && strings.ToLower(ss[i]) == strings.ToLower(alias)) {
			return true
		}
	}
	return false
}
