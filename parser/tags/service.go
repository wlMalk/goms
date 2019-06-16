package tags

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

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

func ServiceGenerateTag(service *types.Service, tag string) error {
	options := strings.SplitS(tag, ",")
	return serviceGenerateOptions.only(service, true, "generate", options...)
}

func ServiceGenerateAllTag(service *types.Service, tag string) error {
	options := strings.SplitS(tag, ",")
	return serviceGenerateOptions.allBut(service, true, "generate-all", options...)
}

func ServiceTransportsTag(service *types.Service, tag string) error {
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

func ServiceHTTPUriPrefixTag(service *types.Service, tag string) error {
	service.Options.HTTP.URIPrefix = tag
	return nil
}

func ServiceNameTag(service *types.Service, tag string) error {
	tag = strs.TrimSpace(tag)
	if len(tag) == 0 {
		return fmt.Errorf("invalid name '%s' for name tag in '%s' service", tag, service.Name)
	}
	service.Alias = tag
	return nil
}

func ServiceMetricsTag(service *types.Service, tag string) error {
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
