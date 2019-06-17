package generators

import (
	"strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceMainFunc(file file.File, service types.Service) {
	file.AddImport("", "io")
	file.AddImport("", "os")
	if helpers.IsServerEnabled(service) {
		file.AddImport("", service.ImportPath, "/cmd/start")
	}
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		file.AddImport("", "github.com/go-kit/kit/log")
	}
	if helpers.IsMetricsEnabled(service) {
		file.AddImport("", "github.com/go-kit/kit/metrics")
	}
	if helpers.IsTracingEnabled(service) {
		file.AddImport("opentracinggo", "github.com/opentracing/opentracing-go")
	}
	if service.Options.Generate.ProtoBuf && (helpers.IsGRPCServerEnabled(service) || helpers.IsGRPCClientEnabled(service)) {
		file.Pf("//go: protoc --go_out=plugins=grpc:%s --proto_path=%s proto/service.goms.proto", strings.TrimSuffix(file.Base(), service.ImportPath), file.Base())
		file.P("")
	}
	file.Pf("func main() {")
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		file.Pf("logger := InitLogger(os.Stderr)")
	}
	if helpers.IsTracingEnabled(service) {
		file.Pf("tracer := InitTracer()")
	}
	if helpers.IsFrequencyMetricEnabled(service) {
		file.Pf("frequencyMetric := InitRequestFrequencyMetric()")
	}
	if helpers.IsLatencyMetricEnabled(service) {
		file.Pf("latencyMetric := InitRequestLatencyMetric()")
	}
	if helpers.IsCounterMetricEnabled(service) {
		file.Pf("counterMetric := InitRequestCounterMetric()")
	}
	if helpers.IsServerEnabled(service) {
		file.Pf("start.Start(")
		if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
			file.Pf("logger,")
		}
		if helpers.IsTracingEnabled(service) {
			file.Pf("tracer,")
		}
		if helpers.IsFrequencyMetricEnabled(service) {
			file.Pf("frequencyMetric,")
		}
		if helpers.IsLatencyMetricEnabled(service) {
			file.Pf("latencyMetric,")
		}
		if helpers.IsCounterMetricEnabled(service) {
			file.Pf("counterMetric,")
		}
		file.Pf(")")
	}
	file.Pf("}")
	file.Pf("")
}

func ServiceMainInitLoggerFunc(file file.File, service types.Service) {
	file.Pf("func InitLogger(writer io.Writer) log.Logger {")
	file.Pf("logger := log.NewJSONLogger(writer)")
	file.Pf("logger = log.With(logger, \"@timestamp\", log.DefaultTimestampUTC)")
	file.Pf("logger = log.With(logger, \"caller\", log.DefaultCaller)")
	file.Pf("return logger")
	file.Pf("}")
	file.Pf("")
}

func ServiceMainInitTracerFunc(file file.File, service types.Service) {
	file.Pf("func InitTracer() opentracinggo.Tracer {")
	file.Pf("// TODO: Initialize tracer")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func ServiceMainInitCounterFunc(file file.File, service types.Service) {
	file.Pf("func InitRequestCounterMetric() metrics.Counter {")
	file.Pf("// TODO: Initialize counterMetric")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func ServiceMainInitLatencyFunc(file file.File, service types.Service) {
	file.Pf("func InitRequestLatencyMetric() metrics.Histogram {")
	file.Pf("// TODO: Initialize latencyMetric")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func ServiceMainInitFrequencyFunc(file file.File, service types.Service) {
	file.Pf("func InitRequestFrequencyMetric() metrics.Gauge {")
	file.Pf("// TODO: Initialize frequencyMetric")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}
