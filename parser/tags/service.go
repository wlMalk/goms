package tags

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/constants"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceTransportsTag(service *types.Service, tag string) error {
	transports := strings.SplitS(tag, ",")
	service.Generate.Remove(
		constants.ServiceGenerateHTTPServerFlag,
		constants.ServiceGenerateHTTPClientFlag,
		constants.ServiceGenerateGRPCServerFlag,
		constants.ServiceGenerateGRPCClientFlag,
	)
	for _, i := range transports {
		switch strs.ToUpper(i) {
		case "HTTP":
			service.Generate.Add(
				constants.ServiceGenerateHTTPServerFlag,
				constants.ServiceGenerateHTTPClientFlag,
			)
		case "GRPC":
			service.Generate.Add(
				constants.ServiceGenerateGRPCServerFlag,
				constants.ServiceGenerateGRPCClientFlag,
			)
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
	service.Generate.Remove(
		constants.ServiceGenerateFrequencyMetricFlag,
		constants.ServiceGenerateLatencyMetricFlag,
		constants.ServiceGenerateCounterMetricFlag,
	)
	for _, i := range transports {
		switch strs.ToLower(i) {
		case "frequency":
			service.Generate.Add(constants.ServiceGenerateFrequencyMetricFlag)
		case "latency":
			service.Generate.Add(constants.ServiceGenerateLatencyMetricFlag)
		case "counter":
			service.Generate.Add(constants.ServiceGenerateCounterMetricFlag)
		default:
			return fmt.Errorf("invalid value '%s' for metrics service tag in '%s' service", i, service.Name)
		}
	}
	return nil
}
