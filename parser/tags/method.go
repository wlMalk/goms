package tags

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

type methodGenerateOptionsHandler map[string]func(method *types.Method, v bool)

var methodGenerateOptions = methodGenerateOptionsHandler{
	"circuit-breaking": func(method *types.Method, v bool) {
		method.Options.Generate.CircuitBreaking = v
	}, "rate-limiting": func(method *types.Method, v bool) {
		method.Options.Generate.RateLimiting = v
	}, "recovering": func(method *types.Method, v bool) {
		method.Options.Generate.Recovering = v
	}, "caching": func(method *types.Method, v bool) {
		method.Options.Generate.Caching = v
	}, "logging": func(method *types.Method, v bool) {
		method.Options.Generate.Logging = v
	}, "tracing": func(method *types.Method, v bool) {
		method.Options.Generate.Tracing = v
	}, "metrics": func(method *types.Method, v bool) {
		method.Options.Generate.FrequencyMetric = v
		method.Options.Generate.LatencyMetric = v
		method.Options.Generate.CounterMetric = v
	}, "validators": func(method *types.Method, v bool) {
		method.Options.Generate.Validators = v
	}, "validating": func(method *types.Method, v bool) {
		method.Options.Generate.Validating = v
	}, "middleware": func(method *types.Method, v bool) {
		method.Options.Generate.Middleware = v
	}, "method-stubs": func(method *types.Method, v bool) {
		method.Options.Generate.MethodStubs = v
	}, "grpc-server": func(method *types.Method, v bool) {
		method.Options.Generate.GRPCServer = v
	}, "grpc-client": func(method *types.Method, v bool) {
		method.Options.Generate.GRPCClient = v
	}, "http-server": func(method *types.Method, v bool) {
		method.Options.Generate.HTTPServer = v
	}, "http-client": func(method *types.Method, v bool) {
		method.Options.Generate.HTTPClient = v
	}, "grpc": func(method *types.Method, v bool) {
		method.Options.Generate.GRPCServer = v
		method.Options.Generate.GRPCClient = v
	}, "http": func(method *types.Method, v bool) {
		method.Options.Generate.HTTPServer = v
		method.Options.Generate.HTTPClient = v
	},
}

func (m methodGenerateOptionsHandler) all(method *types.Method, v bool, tagName string) {
	for _, f := range m {
		f(method, v)
	}
}

func (m methodGenerateOptionsHandler) allBut(method *types.Method, v bool, tagName string, options ...string) error {
	m.all(method, v, tagName)
	return m.only(method, !v, tagName, options...)
}

func (m methodGenerateOptionsHandler) only(method *types.Method, v bool, tagName string, options ...string) error {
	for _, option := range options {
		f, ok := m[strs.ToLower(option)]
		if !ok {
			return fmt.Errorf("invalid value '%s' for %s method tag in '%s' method", option, tagName, method.Name)
		}
		f(method, v)
	}
	return nil
}

func MethodTransportsTag(method *types.Method, tag string) error {
	transports := strings.SplitS(tag, ",")
	method.Options.Generate.HTTPServer = false
	method.Options.Generate.HTTPClient = false
	method.Options.Generate.GRPCServer = false
	method.Options.Generate.GRPCClient = false
	for _, i := range transports {
		switch strs.ToUpper(i) {
		case "HTTP":
			method.Options.Generate.HTTPServer = true
			method.Options.Generate.HTTPClient = true
		case "GRPC":
			method.Options.Generate.GRPCServer = true
			method.Options.Generate.GRPCClient = true
		default:
			return fmt.Errorf("invalid value '%s' for transports method tag in '%s' method", i, method.Name)
		}
	}
	return nil
}

func MethodMetricsTag(method *types.Method, tag string) error {
	transports := strings.SplitS(tag, ",")
	method.Options.Generate.FrequencyMetric = false
	method.Options.Generate.LatencyMetric = false
	method.Options.Generate.CounterMetric = false
	for _, i := range transports {
		switch strs.ToLower(i) {
		case "frequency":
			method.Options.Generate.FrequencyMetric = true
		case "latency":
			method.Options.Generate.LatencyMetric = true
		case "counter":
			method.Options.Generate.CounterMetric = true
		default:
			return fmt.Errorf("invalid value '%s' for metrics method tag in '%s' method", i, method.Name)
		}
	}
	return nil
}

func MethodHTTPMethodTag(method *types.Method, tag string) error {
	httpMethod := strs.ToUpper(tag)
	if httpMethod != "POST" && httpMethod != "GET" && httpMethod != "PUT" && httpMethod != "DELETE" && httpMethod != "OPTIONS" && httpMethod != "HEAD" {
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

func MethodEnableTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.only(method, true, "enable", options...)
}

func MethodDisableTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.only(method, false, "disable", options...)
}

func MethodEnableAllTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.allBut(method, true, "enable-all", options...)
}

func MethodDisableAllTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.allBut(method, false, "disable-all", options...)
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
