package tags

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/constants"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func MethodTransportsTag(method *types.Method, tag string) error {
	transports := strings.SplitS(tag, ",")
	method.Generate.Remove(
		constants.MethodGenerateHTTPServerFlag,
		constants.MethodGenerateHTTPClientFlag,
		constants.MethodGenerateGRPCServerFlag,
		constants.MethodGenerateGRPCClientFlag,
	)
	for _, i := range transports {
		switch strs.ToUpper(i) {
		case "HTTP":
			method.Generate.Add(
				constants.MethodGenerateHTTPServerFlag,
				constants.MethodGenerateHTTPClientFlag,
			)
		case "GRPC":
			method.Generate.Add(
				constants.MethodGenerateGRPCServerFlag,
				constants.MethodGenerateGRPCClientFlag,
			)
		default:
			return fmt.Errorf("invalid value '%s' for transports method tag in '%s' method", i, method.Name)
		}
	}
	return nil
}

func MethodMetricsTag(method *types.Method, tag string) error {
	transports := strings.SplitS(tag, ",")
	method.Generate.Remove(
		constants.MethodGenerateFrequencyMetricFlag,
		constants.MethodGenerateLatencyMetricFlag,
		constants.MethodGenerateCounterMetricFlag,
	)
	for _, i := range transports {
		switch strs.ToLower(i) {
		case "frequency":
			method.Generate.Add(constants.MethodGenerateFrequencyMetricFlag)
		case "latency":
			method.Generate.Add(constants.MethodGenerateLatencyMetricFlag)
		case "counter":
			method.Generate.Add(constants.MethodGenerateCounterMetricFlag)
		default:
			return fmt.Errorf("invalid value '%s' for metrics method tag in '%s' method", i, method.Name)
		}
	}
	return nil
}

func MethodHTTPMethodTag(method *types.Method, tag string) error {
	httpMethod := strs.ToUpper(tag)
	if httpMethod != "POST" &&
		httpMethod != "PATCH" &&
		httpMethod != "GET" &&
		httpMethod != "PUT" &&
		httpMethod != "DELETE" &&
		httpMethod != "OPTIONS" &&
		httpMethod != "HEAD" {
		return fmt.Errorf("invalid http-method value '%s'", tag)
	}
	method.Options.HTTP.Method = httpMethod
	return nil
}

func MethodHTTPUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.URI = tag
	return nil
}

func MethodHTTPAbsUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.AbsURI = tag
	return nil
}

func MethodLogsIgnoreTag(method *types.Method, tag string) error {
	params := strings.SplitS(tag, ",")
paramsLoop:
	for _, p := range params {
		param := strs.ToLower(p)
		if contains(method.Options.Logging.IgnoredArguments, param) || contains(method.Options.Logging.IgnoredResults, param) {
			continue
		}
		if param == "err" {
			method.Options.Logging.IgnoreError = true
			continue
		}
		for _, arg := range method.Arguments {
			if strs.ToLower(arg.Name) == param || (len(arg.Alias) > 0 && strs.ToLower(arg.Alias) == param) {
				method.Options.Logging.IgnoredArguments = append(method.Options.Logging.IgnoredArguments, param)
				continue paramsLoop
			}
		}
		for _, result := range method.Results {
			if strs.ToLower(result.Name) == param || (len(result.Alias) > 0 && strs.ToLower(result.Alias) == param) {
				method.Options.Logging.IgnoredResults = append(method.Options.Logging.IgnoredResults, param)
				continue paramsLoop
			}
		}
		return fmt.Errorf("invalid name '%s' given to logs-ignore method tag in '%s' method", p, method.Name)
	}
	return nil
}

func MethodLogsLenTag(method *types.Method, tag string) error {
	params := strings.SplitS(tag, ",")
paramsLoop:
	for _, p := range params {
		param := strs.ToLower(p)
		if contains(method.Options.Logging.LenArguments, param) || contains(method.Options.Logging.LenResults, param) {
			continue
		}
		for _, arg := range method.Arguments {
			if (arg.Type.IsMap || arg.Type.IsVariadic || arg.Type.IsSlice || arg.Type.IsBytes) &&
				(strs.ToLower(arg.Name) == param || (len(arg.Alias) > 0 && strs.ToLower(arg.Alias) == param)) {
				method.Options.Logging.LenArguments = append(method.Options.Logging.LenArguments, param)
				continue paramsLoop
			}
		}
		for _, result := range method.Results {
			if (result.Type.IsMap || result.Type.IsVariadic || result.Type.IsSlice || result.Type.IsBytes) &&
				(strs.ToLower(result.Name) == param || (len(result.Alias) > 0 && strs.ToLower(result.Alias) == param)) {
				method.Options.Logging.LenResults = append(method.Options.Logging.LenResults, param)
				continue paramsLoop
			}
		}
		return fmt.Errorf("invalid name '%s' given to logs-len method tag in '%s' method", p, method.Name)
	}
	return nil
}

func MethodAliasTag(method *types.Method, tag string) error {
	params := strings.SplitS(tag, ",")
	if len(params) != 2 || strs.TrimSpace(params[0]) == "" || strs.TrimSpace(params[1]) == "" {
		return fmt.Errorf("invalid params '%s' for alias tag in '%s' method", tag, method.Name)
	}
	for _, arg := range method.Arguments {
		if strings.ToUpperFirst(arg.Name) == strings.ToUpperFirst(params[0]) {
			arg.Alias = strs.TrimSpace(params[1])
			return nil
		}
	}
	for _, res := range method.Results {
		if strings.ToUpperFirst(res.Name) == strings.ToUpperFirst(params[0]) {
			res.Alias = strs.TrimSpace(params[1])
			return nil
		}
	}
	return fmt.Errorf("invalid name '%s' for alias tag in '%s' method", params[0], method.Name)
}

func MethodNameTag(method *types.Method, tag string) error {
	tag = strs.TrimSpace(tag)
	if len(tag) == 0 {
		return fmt.Errorf("invalid name '%s' for name tag in '%s' method", tag, method.Name)
	}
	method.Alias = tag
	return nil
}

func MethodValidateTag(method *types.Method, tag string) error {
	return nil
}
