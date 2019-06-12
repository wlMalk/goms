package parser

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

var serviceTags = []string{
	"generate-all",
	"generate",
	"transports",
	"metrics",
	"http-URI-prefix",
}

var methodTags = []string{
	"transports",
	"metrics",
	"http-method",
	"http-URI",
	"http-abs-URI",
	"params",
	"logs-ignore",
	"logs-len",
	"disable",
	"enable",
	"alias",
}

var paramTags = []string{
	"http-origin",
}

var serviceTagsParsers = map[string]func(service *types.Service, tag string) error{
	"generate-all":    parseServiceGenerateAllTag,
	"generate":        parseServiceGenerateTag,
	"transports":      parseServiceTransportsTag,
	"metrics":         parseServiceMetricsTag,
	"http-URI-prefix": parseServiceHttpUriPrefixTag,
}

var methodTagsParsers = map[string]func(method *types.Method, tag string) error{
	"transports":   parseMethodTransportsTag,
	"metrics":      parseMethodMetricsTag,
	"http-method":  parseMethodHttpMethodTag,
	"http-URI":     parseMethodHttpUriTag,
	"http-abs-URI": parseMethodHttpAbsUriTag,
	"params":       parseMethodParamsTag,
	"logs-ignore":  parseMethodLogsIgnoreTag,
	"logs-len":     parseMethodLogsLenTag,
	"disable":      parseMethodDisableTag,
	"enable":       parseMethodEnableTag,
	"alias":        parseMethodAliasTag,
}

var paramTagsParsers = map[string]func(argument *types.Argument, tag string) error{
	"http-origin": parseParamHttpOriginTag,
}

func limitLineLength(str string, length int) []string {
	words := strs.Fields(str)
	var lines []string
	charCount := 0
	var line []string
	for i := 0; i < len(words); {
		count := charCount + len(words[i]) + len(line) - 1
		if count <= length {
			line = append(line, words[i])
			charCount += len(words[i])
			i++
		}
		if count > length || i == len(words) {
			lines = append(lines, strs.Join(line, " "))
			line = []string{}
			charCount = 0
		}
	}
	return lines
}

func cleanComments(comments []string) (tags []string, docs []string) {
	for i := range comments {
		comments[i] = strs.TrimSpace(strs.TrimPrefix(strs.TrimSpace(comments[i]), "//"))
	}
	comments = strings.SplitS(strs.Join(comments, " "), " ")
	for i := range comments {
		comments[i] = strs.TrimSpace(comments[i])
		if comments[i][0] == '@' {
			tags = append(tags, comments[i])
		} else {
			docs = append(docs, comments[i])
		}
	}
	return strings.SplitS(strs.Join(tags, " "), " "), limitLineLength(strs.Join(docs, " "), 80)
}

func parseServiceTags(service *types.Service, tags []string) error {
	for _, tag := range tags {
		found := false
		for _, serviceTag := range serviceTags {
			if strs.HasPrefix(strs.ToLower(tag), "@"+strs.ToLower(serviceTag)) {
				found = true
				if err := serviceTagsParsers[serviceTag](service, cleanTag(strs.TrimPrefix(tag, "@"+serviceTag))); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid service tag \"%s\"", tag)
		}
	}
	return nil
}

func parseMethodTags(method *types.Method, tags []string) error {
	for _, tag := range tags {
		found := false
		for _, methodTag := range methodTags {
			if strs.HasPrefix(tag, "@"+methodTag) {
				found = true
				if err := methodTagsParsers[methodTag](method, cleanTag(strs.TrimPrefix(tag, "@"+methodTag))); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid method tag \"%s\"", tag)
		}
	}
	return nil
}

func parseParamTags(arg *types.Argument, tags []string) error {
	for _, tag := range tags {
		found := false
		for _, paramTag := range paramTags {
			if strs.HasPrefix(tag, "@"+paramTag) {
				found = true
				if err := paramTagsParsers[paramTag](arg, cleanTag(strs.TrimPrefix(tag, "@"+paramTag))); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid param tag \"%s\"", tag)
		}
	}
	return nil
}

func cleanTag(tag string) string {
	return strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(strs.TrimSpace(tag), "("), ")"))
}

func parseServiceGenerateTag(service *types.Service, tag string) error {
	generated := strings.SplitS(tag, ",")
	for _, i := range generated {
		switch strs.ToLower(i) {
		case "logger":
			service.Options.Generate.Logger = true
		case "circuit-breaking":
			service.Options.Generate.CircuitBreaking = true
		case "rate-limiting":
			service.Options.Generate.RateLimiting = true
		case "recovering":
			service.Options.Generate.Recovering = true
		case "caching":
			service.Options.Generate.Caching = true
		case "logging":
			service.Options.Generate.Logging = true
		case "main":
			service.Options.Generate.Main = true
		case "protobuf":
			service.Options.Generate.ProtoBuf = true
		case "tracing":
			service.Options.Generate.Tracing = true
		case "metrics":
			service.Options.Generate.FrequencyMetric = true
			service.Options.Generate.LatencyMetric = true
			service.Options.Generate.CounterMetric = true
		case "service-discovery":
			service.Options.Generate.ServiceDiscovery = true
		case "validators":
			service.Options.Generate.Validators = true
		case "middleware":
			service.Options.Generate.Middleware = true
		case "method-stubs":
			service.Options.Generate.MethodStubs = true
		case "grpc-server":
			service.Options.Generate.GRPCServer = true
		case "grpc-client":
			service.Options.Generate.GRPCClient = true
		case "http-server":
			service.Options.Generate.HTTPServer = true
		case "http-client":
			service.Options.Generate.HTTPClient = true
		case "grpc":
			service.Options.Generate.GRPCServer = true
			service.Options.Generate.GRPCClient = true
		case "http":
			service.Options.Generate.HTTPServer = true
			service.Options.Generate.HTTPClient = true
		default:
			return fmt.Errorf("invalid value '%s' for generate service tag in '%s' service", i, service.Name)
		}
	}
	return nil
}

func parseServiceGenerateAllTag(service *types.Service, tag string) error {
	ignored := strings.SplitS(tag, ",")
	service.Options.Generate.Logger = true
	service.Options.Generate.CircuitBreaking = true
	service.Options.Generate.RateLimiting = true
	service.Options.Generate.Recovering = true
	service.Options.Generate.Caching = true
	service.Options.Generate.Logging = true
	service.Options.Generate.Main = true
	service.Options.Generate.ProtoBuf = true
	service.Options.Generate.Tracing = true
	service.Options.Generate.FrequencyMetric = true
	service.Options.Generate.LatencyMetric = true
	service.Options.Generate.CounterMetric = true
	service.Options.Generate.ServiceDiscovery = true
	service.Options.Generate.Validators = true
	service.Options.Generate.Middleware = true
	service.Options.Generate.MethodStubs = true
	service.Options.Generate.GRPCServer = true
	service.Options.Generate.GRPCClient = true
	service.Options.Generate.HTTPServer = true
	service.Options.Generate.HTTPClient = true

	for _, i := range ignored {
		switch strs.ToLower(i) {
		case "logger":
			service.Options.Generate.Logger = false
		case "circuit-breaking":
			service.Options.Generate.CircuitBreaking = false
		case "rate-limiting":
			service.Options.Generate.RateLimiting = false
		case "recovering":
			service.Options.Generate.Recovering = false
		case "caching":
			service.Options.Generate.Caching = false
		case "logging":
			service.Options.Generate.Logging = false
		case "main":
			service.Options.Generate.Main = false
		case "protobuf":
			service.Options.Generate.ProtoBuf = false
		case "tracing":
			service.Options.Generate.Tracing = false
		case "metrics":
			service.Options.Generate.FrequencyMetric = false
			service.Options.Generate.LatencyMetric = false
			service.Options.Generate.CounterMetric = false
		case "service-discovery":
			service.Options.Generate.ServiceDiscovery = false
		case "validators":
			service.Options.Generate.Validators = false
		case "middleware":
			service.Options.Generate.Middleware = false
		case "method-stubs":
			service.Options.Generate.MethodStubs = false
		case "grpc-server":
			service.Options.Generate.GRPCServer = false
		case "grpc-client":
			service.Options.Generate.GRPCClient = false
		case "http-server":
			service.Options.Generate.HTTPServer = false
		case "http-client":
			service.Options.Generate.HTTPClient = false
		case "grpc":
			service.Options.Generate.GRPCServer = false
			service.Options.Generate.GRPCClient = false
		case "http":
			service.Options.Generate.HTTPServer = false
			service.Options.Generate.HTTPClient = false
		default:
			return fmt.Errorf("invalid value '%s' for generate-all service tag in '%s' service", i, service.Name)
		}
	}
	return nil
}

func parseServiceTransportsTag(service *types.Service, tag string) error {
	transports := strings.SplitS(tag, ",")
	service.Options.Generate.HTTPServer = false
	service.Options.Generate.HTTPClient = false
	service.Options.Generate.GRPCServer = false
	service.Options.Generate.GRPCClient = false
	for _, i := range transports {
		switch strs.ToUpper(i) {
		case "HTTP":
			service.Options.Generate.HTTPServer = true
			service.Options.Generate.HTTPClient = true
		case "GRPC":
			service.Options.Generate.GRPCServer = true
			service.Options.Generate.GRPCClient = true
		default:
			return fmt.Errorf("invalid value '%s' for transports service tag in '%s' service", i, service.Name)
		}
	}
	return nil
}

func parseServiceHttpUriPrefixTag(service *types.Service, tag string) error {
	service.Options.HTTP.URIPrefix = tag
	return nil
}

func parseMethodTransportsTag(method *types.Method, tag string) error {
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

func parseMethodMetricsTag(method *types.Method, tag string) error {
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

func parseServiceMetricsTag(service *types.Service, tag string) error {
	transports := strings.SplitS(tag, ",")
	service.Options.Generate.FrequencyMetric = false
	service.Options.Generate.LatencyMetric = false
	service.Options.Generate.CounterMetric = false
	for _, i := range transports {
		switch strs.ToLower(i) {
		case "frequency":
			service.Options.Generate.FrequencyMetric = true
		case "latency":
			service.Options.Generate.LatencyMetric = true
		case "counter":
			service.Options.Generate.CounterMetric = true
		default:
			return fmt.Errorf("invalid value '%s' for metrics service tag in '%s' service", i, service.Name)
		}
	}
	return nil
}

func parseMethodHttpMethodTag(method *types.Method, tag string) error {
	httpMethod := strs.ToUpper(tag)
	if httpMethod != "POST" && httpMethod != "GET" && httpMethod != "PUT" && httpMethod != "DELETE" && httpMethod != "OPTIONS" && httpMethod != "HEAD" {
		return fmt.Errorf("invalid http-method value '%s'", tag)
	}
	method.Options.HTTP.Method = httpMethod
	return nil
}

func parseMethodHttpUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.URI = tag
	return nil
}

func parseMethodHttpAbsUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.AbsURI = tag
	return nil
}

func errInvalidTagFormat(origin string, tagName string, tag string) error {
	return fmt.Errorf("invalid %s tag '%s': @%s(%s)", origin, tagName, tagName, tag)
}

func parseMethodParamsTag(method *types.Method, tag string) error {
	params := strings.SplitS(tag, ",")
	if len(params) != 2 {
		return errInvalidTagFormat("method", "params", tag)
	}
	argNames := strings.SplitS(strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(params[0], "["), "]")), ",")
	args := filterArguments(method.Arguments, func(arg *types.Argument) bool {
		return contains(argNames, strings.ToLowerFirst(arg.Name))
	})
	tags := strings.SplitS(strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(params[1], "("), ")")), ",")
	for _, arg := range args {
		err := parseParamTags(arg, tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseMethodLogsIgnoreTag(method *types.Method, tag string) error {
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

func parseMethodLogsLenTag(method *types.Method, tag string) error {
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

func parseMethodDisableTag(method *types.Method, tag string) error {
	disabled := strings.SplitS(tag, ",")
	for _, i := range disabled {
		switch strs.ToLower(i) {
		case "circuit-breaking":
			method.Options.Generate.CircuitBreaking = false
		case "rate-limiting":
			method.Options.Generate.RateLimiting = false
		case "recovering":
			method.Options.Generate.Recovering = false
		case "caching":
			method.Options.Generate.Caching = false
		case "logging":
			method.Options.Generate.Logging = false
		case "tracing":
			method.Options.Generate.Tracing = false
		case "metrics":
			method.Options.Generate.FrequencyMetric = false
			method.Options.Generate.LatencyMetric = false
			method.Options.Generate.CounterMetric = false
		case "validators":
			method.Options.Generate.Validator = false
		case "middleware":
			method.Options.Generate.Middleware = false
		case "method-stubs":
			method.Options.Generate.MethodStubs = false
		case "grpc-server":
			method.Options.Generate.GRPCServer = false
		case "grpc-client":
			method.Options.Generate.GRPCClient = false
		case "http-server":
			method.Options.Generate.HTTPServer = false
		case "http-client":
			method.Options.Generate.HTTPClient = false
		case "grpc":
			method.Options.Generate.GRPCServer = false
			method.Options.Generate.GRPCClient = false
		case "http":
			method.Options.Generate.HTTPServer = false
			method.Options.Generate.HTTPClient = false
		default:
			return fmt.Errorf("invalid value '%s' for disable method tag in '%s' method", i, method.Name)
		}
	}
	return nil
}

func parseMethodEnableTag(method *types.Method, tag string) error {
	enabled := strings.SplitS(tag, ",")
	for _, i := range enabled {
		switch strs.ToLower(i) {
		case "circuit-breaking":
			method.Options.Generate.CircuitBreaking = true
		case "rate-limiting":
			method.Options.Generate.RateLimiting = true
		case "recovering":
			method.Options.Generate.Recovering = true
		case "caching":
			method.Options.Generate.Caching = true
		case "logging":
			method.Options.Generate.Logging = true
		case "tracing":
			method.Options.Generate.Tracing = true
		case "metrics":
			method.Options.Generate.FrequencyMetric = true
			method.Options.Generate.LatencyMetric = true
			method.Options.Generate.CounterMetric = true
		case "validators":
			method.Options.Generate.Validator = true
		case "middleware":
			method.Options.Generate.Middleware = true
		case "method-stubs":
			method.Options.Generate.MethodStubs = true
		case "grpc-server":
			method.Options.Generate.GRPCServer = true
		case "grpc-client":
			method.Options.Generate.GRPCClient = true
		case "http-server":
			method.Options.Generate.HTTPServer = true
		case "http-client":
			method.Options.Generate.HTTPClient = true
		case "grpc":
			method.Options.Generate.GRPCServer = true
			method.Options.Generate.GRPCClient = true
		case "http":
			method.Options.Generate.HTTPServer = true
			method.Options.Generate.HTTPClient = true

		default:
			return fmt.Errorf("invalid value '%s' for enable method tag in '%s' method", i, method.Name)
		}
	}
	return nil
}

func parseMethodAliasTag(method *types.Method, tag string) error {
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

func parseMethodValidateTag(method *types.Method, tag string) error {
	return nil
}

func parseParamHttpOriginTag(arg *types.Argument, tag string) error {
	origin := strs.ToUpper(tag)
	if origin != "BODY" && origin != "HEADER" && origin != "QUERY" && origin != "PATH" {
		return fmt.Errorf("invalid http-origin value '%s'", tag)
	}
	arg.Options.HTTP.Origin = origin
	return nil
}

func filterArguments(args []*types.Argument, f func(*types.Argument) bool) (filtered []*types.Argument) {
	for _, arg := range args {
		if f(arg) {
			filtered = append(filtered, arg)
		}
	}
	return
}

func filterFields(fields []*types.Field, f func(*types.Field) bool) (filtered []*types.Field) {
	for _, field := range fields {
		if f(field) {
			filtered = append(filtered, field)
		}
	}
	return
}
