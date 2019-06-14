package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateServiceMainFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, false, false)
	file.Pkg = "main"
	generateServiceMainFunc(file, service)
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		generateServiceMainInitLoggerFunc(file, service)
	}
	if helpers.IsTracingEnabled(service) {
		generateServiceMainInitTracerFunc(file, service)
	}
	if helpers.IsFrequencyMetricEnabled(service) {
		generateServiceMainInitFrequencyFunc(file, service)
	}
	if helpers.IsLatencyMetricEnabled(service) {
		generateServiceMainInitLatencyFunc(file, service)
	}
	if helpers.IsCounterMetricEnabled(service) {
		generateServiceMainInitCounterFunc(file, service)
	}
	return file
}

func generateServiceMainFunc(file *files.GoFile, service *types.Service) {
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

func generateServiceMainInitLoggerFunc(file *files.GoFile, service *types.Service) {
	file.Pf("func InitLogger(writer io.Writer) log.Logger {")
	file.Pf("logger := log.NewJSONLogger(writer)")
	file.Pf("logger = log.With(logger, \"@timestamp\", log.DefaultTimestampUTC)")
	file.Pf("logger = log.With(logger, \"caller\", log.DefaultCaller)")
	file.Pf("return logger")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainInitTracerFunc(file *files.GoFile, service *types.Service) {
	file.Pf("func InitTracer() opentracinggo.Tracer {")
	file.Pf("// TODO: Initialize tracer")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainInitCounterFunc(file *files.GoFile, service *types.Service) {
	file.Pf("func InitRequestCounterMetric() metrics.Counter {")
	file.Pf("// TODO: Initialize counterMetric")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainInitLatencyFunc(file *files.GoFile, service *types.Service) {
	file.Pf("func InitRequestLatencyMetric() metrics.Histogram {")
	file.Pf("// TODO: Initialize latencyMetric")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainInitFrequencyFunc(file *files.GoFile, service *types.Service) {
	file.Pf("func InitRequestFrequencyMetric() metrics.Gauge {")
	file.Pf("// TODO: Initialize frequencyMetric")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}
