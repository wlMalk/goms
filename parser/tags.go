package parser

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

var serviceTags = []string{
	"name",
	"generate-all",
	"generate",
	"transports",
	"metrics",
	"http-URI-prefix",
}

var methodTags = []string{
	"name",
	"transports",
	"metrics",
	"http-method",
	"http-URI",
	"http-abs-URI",
	"params",
	"logs-ignore",
	"logs-len",
	"disable-all",
	"enable-all",
	"disable",
	"enable",
	"alias",
}

var paramTags = []string{
	"http-origin",
}

var serviceTagsParsers = map[string]func(service *types.Service, tag string) error{
	"name":            parseServiceNameTag,
	"generate-all":    parseServiceGenerateAllTag,
	"generate":        parseServiceGenerateTag,
	"transports":      parseServiceTransportsTag,
	"metrics":         parseServiceMetricsTag,
	"http-URI-prefix": parseServiceHTTPUriPrefixTag,
}

var methodTagsParsers = map[string]func(method *types.Method, tag string) error{
	"name":         parseMethodNameTag,
	"transports":   parseMethodTransportsTag,
	"metrics":      parseMethodMetricsTag,
	"http-method":  parseMethodHTTPMethodTag,
	"http-URI":     parseMethodHTTPUriTag,
	"http-abs-URI": parseMethodHTTPAbsUriTag,
	"params":       parseMethodParamsTag,
	"logs-ignore":  parseMethodLogsIgnoreTag,
	"logs-len":     parseMethodLogsLenTag,
	"disable":      parseMethodDisableTag,
	"enable":       parseMethodEnableTag,
	"disable-all":  parseMethodDisableAllTag,
	"enable-all":   parseMethodEnableAllTag,
	"alias":        parseMethodAliasTag,
}

var paramTagsParsers = map[string]func(argument *types.Argument, tag string) error{
	"http-origin": parseParamHTTPOriginTag,
}

type serviceGenerateOptionsHandler map[string]func(service *types.Service, v bool)

var serviceGenerateOptions = serviceGenerateOptionsHandler{
	"logger": func(service *types.Service, v bool) {
		service.Options.Generate.Logger = v
	}, "circuit-breaking": func(service *types.Service, v bool) {
		service.Options.Generate.CircuitBreaking = v
	}, "rate-limiting": func(service *types.Service, v bool) {
		service.Options.Generate.RateLimiting = v
	}, "recovering": func(service *types.Service, v bool) {
		service.Options.Generate.Recovering = v
	}, "caching": func(service *types.Service, v bool) {
		service.Options.Generate.Caching = v
	}, "logging": func(service *types.Service, v bool) {
		service.Options.Generate.Logging = v
	}, "main": func(service *types.Service, v bool) {
		service.Options.Generate.Main = v
	}, "protobuf": func(service *types.Service, v bool) {
		service.Options.Generate.ProtoBuf = v
	}, "tracing": func(service *types.Service, v bool) {
		service.Options.Generate.Tracing = v
	}, "metrics": func(service *types.Service, v bool) {
		service.Options.Generate.FrequencyMetric = v
		service.Options.Generate.LatencyMetric = v
		service.Options.Generate.CounterMetric = v
	}, "service-discovery": func(service *types.Service, v bool) {
		service.Options.Generate.ServiceDiscovery = v
	}, "validators": func(service *types.Service, v bool) {
		service.Options.Generate.Validators = v
	}, "validating": func(service *types.Service, v bool) {
		service.Options.Generate.Validating = v
	}, "middleware": func(service *types.Service, v bool) {
		service.Options.Generate.Middleware = v
	}, "method-stubs": func(service *types.Service, v bool) {
		service.Options.Generate.MethodStubs = v
	}, "grpc-server": func(service *types.Service, v bool) {
		service.Options.Generate.GRPCServer = v
	}, "grpc-client": func(service *types.Service, v bool) {
		service.Options.Generate.GRPCClient = v
	}, "http-server": func(service *types.Service, v bool) {
		service.Options.Generate.HTTPServer = v
	}, "http-client": func(service *types.Service, v bool) {
		service.Options.Generate.HTTPClient = v
	}, "grpc": func(service *types.Service, v bool) {
		service.Options.Generate.GRPCServer = v
		service.Options.Generate.GRPCClient = v
	}, "http": func(service *types.Service, v bool) {
		service.Options.Generate.HTTPServer = v
		service.Options.Generate.HTTPClient = v
	}, "dockerfile": func(service *types.Service, v bool) {
		service.Options.Generate.Dockerfile = v
	},
}

func (m serviceGenerateOptionsHandler) all(service *types.Service, v bool, tagName string) {
	for _, f := range m {
		f(service, v)
	}
}

func (m serviceGenerateOptionsHandler) allBut(service *types.Service, v bool, tagName string, options ...string) error {
	m.all(service, v, tagName)
	return m.only(service, !v, tagName, options...)
}

func (m serviceGenerateOptionsHandler) only(service *types.Service, v bool, tagName string, options ...string) error {
	for _, option := range options {
		f, ok := m[strs.ToLower(option)]
		if !ok {
			return fmt.Errorf("invalid value '%s' for %s service tag in '%s' service", option, tagName, service.Name)
		}
		f(service, v)
	}
	return nil
}

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
		comments[i] = strs.TrimSpace(comments[i])
		comments[i] = strs.Replace(comments[i], "\n", "", -1)
		comments[i] = strs.Replace(comments[i], "\t", " ", -1)
		comments[i] = strs.TrimPrefix(strs.TrimPrefix(strs.TrimSuffix(comments[i], "*/"), "/*"), "//")
		comments[i] = strs.TrimSpace(comments[i])
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
	options := strings.SplitS(tag, ",")
	return serviceGenerateOptions.only(service, true, "generate", options...)
}

func parseServiceGenerateAllTag(service *types.Service, tag string) error {
	options := strings.SplitS(tag, ",")
	return serviceGenerateOptions.allBut(service, true, "generate-all", options...)
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

func parseServiceHTTPUriPrefixTag(service *types.Service, tag string) error {
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

func parseMethodHTTPMethodTag(method *types.Method, tag string) error {
	httpMethod := strs.ToUpper(tag)
	if httpMethod != "POST" && httpMethod != "GET" && httpMethod != "PUT" && httpMethod != "DELETE" && httpMethod != "OPTIONS" && httpMethod != "HEAD" {
		return fmt.Errorf("invalid http-method value '%s'", tag)
	}
	method.Options.HTTP.Method = httpMethod
	return nil
}

func parseMethodHTTPUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.URI = tag
	return nil
}

func parseMethodHTTPAbsUriTag(method *types.Method, tag string) error {
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
	if len(argNames) != len(args) {
		return fmt.Errorf("invalid arguments '%s' given to params method tag in '%s' method", strs.Join(argNames, ", "), method.Name)
	}
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

func parseMethodEnableTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.only(method, true, "enable", options...)
}

func parseMethodDisableTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.only(method, false, "disable", options...)
}

func parseMethodEnableAllTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.allBut(method, true, "enable-all", options...)
}

func parseMethodDisableAllTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	return methodGenerateOptions.allBut(method, false, "disable-all", options...)
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

func parseMethodNameTag(method *types.Method, tag string) error {
	tag = strs.TrimSpace(tag)
	if len(tag) == 0 {
		return fmt.Errorf("invalid name '%s' for name tag in '%s' method", tag, method.Name)
	}
	method.Alias = tag
	return nil
}

func parseServiceNameTag(service *types.Service, tag string) error {
	tag = strs.TrimSpace(tag)
	if len(tag) == 0 {
		return fmt.Errorf("invalid name '%s' for name tag in '%s' service", tag, service.Name)
	}
	service.Alias = tag
	return nil
}

func parseMethodValidateTag(method *types.Method, tag string) error {
	return nil
}

func parseParamHTTPOriginTag(arg *types.Argument, tag string) error {
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
