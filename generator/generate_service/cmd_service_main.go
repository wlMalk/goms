package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
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
	generateServiceMainInitEndpointsFunc(file, service)
	generateServiceMainPrepareEndpointsFunc(file, service)
	if helpers.IsServerEnabled(service) {
		generateServiceMainInterruptHandlerFunc(file, service)
	}
	if helpers.IsHTTPServerEnabled(service) {
		generateServiceMainServeHTTPFunc(file, service)
	}
	return file
}

func generateServiceMainFunc(file *files.GoFile, service *types.Service) {
	serviceNameSnake := strings.ToSnakeCase(service.Name)

	file.AddImport("", "context")
	file.AddImport("", "errors")
	file.AddImport("", "fmt")
	file.AddImport("", "io")
	file.AddImport("", "net/http")
	file.AddImport("", "os")
	file.AddImport("", "os/signal")
	file.AddImport("", "syscall")
	file.AddImport("", "time")

	file.AddImport("", "golang.org/x/sync/errgroup")
	file.AddImport("", "golang.org/x/time/rate")

	file.AddImport("", "github.com/go-kit/kit/circuitbreaker")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	file.AddImport("", "github.com/go-kit/kit/log")
	file.AddImport("", "github.com/go-kit/kit/metrics")
	file.AddImport("", "github.com/go-kit/kit/ratelimit")
	file.AddImport("", "github.com/go-kit/kit/tracing/opentracing")
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")

	file.AddImport("", service.ImportPath, "/service/handlers/converters")
	file.AddImport("", service.ImportPath, "/service/middleware")
	file.AddImport("", service.ImportPath, "/service/transport")
	file.AddImport(strings.ToSnakeCase(service.Name)+"_http_server", service.ImportPath, "/service/transport/http/server")
	file.AddImport("", service.ImportPath, "/"+strings.ToSnakeCase(service.Name))

	file.AddImport("goms_middleware", "github.com/wlMalk/goms/goms/middleware")
	file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
	file.AddImport("goms_router", "github.com/wlMalk/goms/goms/transport/http/httprouter")

	file.AddImport("", "github.com/julienschmidt/httprouter")
	file.AddImport("opentracinggo", "github.com/opentracing/opentracing-go")
	file.AddImport("", "github.com/sony/gobreaker")

	file.Pf("func main() {")
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		file.Pf("logger := InitLogger(os.Stderr)")
		file.Pf("")
	}
	if helpers.IsTracingEnabled(service) {
		file.Pf("tracer := InitTracer()")
		file.Pf("")
	}
	if helpers.IsFrequencyMetricEnabled(service) {
		file.Pf("frequencyMetric := InitRequestFrequencyMetric()")
		file.Pf("")
	}
	if helpers.IsLatencyMetricEnabled(service) {
		file.Pf("latencyMetric := InitRequestLatencyMetric()")
		file.Pf("")
	}
	if helpers.IsCounterMetricEnabled(service) {
		file.Pf("counterMetric := InitRequestCounterMetric()")
		file.Pf("")
	}
	if service.Options.Generate.Logger {
		file.Pf("logger.Log(\"message\", \"Hello, I am alive\")")
		file.Pf("defer logger.Log(\"message\", \"goodbye, good luck\")")
		file.Pf("")
	}

	if helpers.IsServerEnabled(service) {
		file.Pf("g, ctx := errgroup.WithContext(context.Background())")
		file.Pf("g.Go(func() error {")
		file.Pf("return InterruptHandler(ctx)")
		file.Pf("})")
		file.Pf("")
	}
	file.Pf("s := %s.New()", serviceNameSnake)
	file.Pf("endpoints := InitEndpoints(s)")
	file.Pf("endpoints = PrepareEndpoints(")
	file.Pf("endpoints,")
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
	if helpers.IsHTTPServerEnabled(service) {
		file.Pf("")
		file.Pf("httpAddr := \":8080\" // TODO: use normal address")
		file.Pf("g.Go(func() error {")
		file.Pf("return ServeHTTP(")
		file.Pf("ctx,")
		file.Pf("&endpoints,")
		file.Pf("httpAddr,")
		if service.Options.Generate.Logger {
			file.Pf("log.With(logger, \"transport\", \"HTTP\"),")
		}
		if helpers.IsTracingEnabled(service) && service.Options.Generate.Logger {
			file.Pf("tracer,")
		}
		file.Pf(")")
		file.Pf("})")
	}
	if helpers.IsServerEnabled(service) {
		file.Pf("")
		if service.Options.Generate.Logger {
			file.Pf("if err := g.Wait(); err != nil {")
			file.Pf("logger.Log(\"error\", err)")
			file.Pf("}")
		} else {
			file.Pf("g.Wait()")
		}
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

func generateServiceMainInitEndpointsFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.Pf("func InitEndpoints(s *%s.%s) transport.%s {", serviceNameSnake, serviceName, serviceName)
	file.Pf("return transport.Endpoints(")
	file.Pf("converters.RequestResponseHandlerToEndpointHandler(")
	file.Pf("converters.RequestHandlerToRequestResponseHandler(")
	file.Pf("converters.HandlerToRequestHandler(s))),")
	file.Pf("s,")
	file.Pf(")")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainPrepareEndpointsFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("func PrepareEndpoints(")
	file.Pf("endpoints transport.%s,", serviceName)
	if helpers.IsTracingEnabled(service) {
		file.Pf("tracer opentracinggo.Tracer,")
	}
	if helpers.IsFrequencyMetricEnabled(service) {
		file.Pf("frequencyMetric metrics.Gauge,")
	}
	if helpers.IsLatencyMetricEnabled(service) {
		file.Pf("latencyMetric metrics.Histogram,")
	}
	if helpers.IsCounterMetricEnabled(service) {
		file.Pf("counterMetric metrics.Counter,")
	}
	file.Pf(") transport.%s {", serviceName)
	file.Pf("")
	if helpers.IsRateLimitingEnabled(service) ||
		helpers.IsCircuitBreakingEnabled(service) ||
		helpers.IsRecoveringEnabled(service) ||
		helpers.IsLoggingEnabled(service) ||
		helpers.IsTracingEnabled(service) ||
		helpers.IsMetricsEnabled(service) {
		file.Pf("endpoints = middleware.ApplyMiddlewareSpecial(endpoints,")
		file.Pf("func(method string) (mw []endpoint.Middleware) {")
		if helpers.IsRateLimitingEnabled(service) {
			file.Pf("mw = append(mw, ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1)))")
		}
		if helpers.IsCircuitBreakingEnabled(service) {
			file.Pf("mw = append(mw, circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})))")
		}
		if helpers.IsRecoveringEnabled(service) {
			file.Pf("mw = append(mw, goms_middleware.RecoveringMiddleware())")
		}
		if helpers.IsLoggingEnabled(service) {
			file.Pf("mw = append(mw, goms_middleware.LoggingMiddleware())")
		}
		if helpers.IsTracingEnabled(service) {
			file.Pf("mw = append(mw, opentracing.TraceServer(tracer, \"%s.\"+method))", serviceName)
		}
		if helpers.IsMetricsEnabled(service) {
			file.Pf("mw = append(mw, goms_middleware.InstrumentingMiddleware(")
			if helpers.IsFrequencyMetricEnabled(service) {
				file.Pf("frequencyMetric.With(\"service\", \"%s\", \"method\", method),", helpers.GetName(serviceName, service.Alias))
			}
			if helpers.IsLatencyMetricEnabled(service) {
				file.Pf("latencyMetric.With(\"service\", \"%s\", \"method\", method),", helpers.GetName(serviceName, service.Alias))
			}
			if helpers.IsCounterMetricEnabled(service) {
				file.Pf("counterMetric.With(\"service\", \"%s\", \"method\", method),", helpers.GetName(serviceName, service.Alias))
			}
			file.Pf("))")
		}

		file.Pf("return")
		file.Pf("},")
		file.Pf(")")
	}
	file.Pf("return endpoints")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainInterruptHandlerFunc(file *files.GoFile, service *types.Service) {
	file.Pf("func InterruptHandler(ctx context.Context) error {")
	file.Pf("interruptHandler := make(chan os.Signal, 1)")
	file.Pf("signal.Notify(interruptHandler, syscall.SIGINT, syscall.SIGTERM)")
	file.Pf("select {")
	file.Pf("case sig := <-interruptHandler:")
	file.Pf("return fmt.Errorf(\"signal received: %%v\", sig.String())")
	file.Pf("case <-ctx.Done():")
	file.Pf("return errors.New(\"signal listener: context canceled\")")
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}

func generateServiceMainServeHTTPFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.Pf("func ServeHTTP(")
	file.Pf("ctx context.Context,")
	file.Pf("endpoints *transport.Strings,")
	file.Pf("addr string,")
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		file.Pf("logger log.Logger,")
	}
	if helpers.IsTracingEnabled(service) && service.Options.Generate.Logger {
		file.Pf("tracer opentracinggo.Tracer,")
	}
	file.Pf(") error {")
	file.Pf("r := httprouter.New()")
	file.Pf("router := goms_router.New(r)")
	file.Pf("")
	file.Pf("server := goms_http.NewServer(router)")
	file.Pf("server.Addr = addr")
	file.Pf("")
	file.Pf("%s_http_server.RegisterSpecial(server, endpoints,", serviceNameSnake)
	file.Pf("func(method string) (opts []kit_http.ServerOption) {")
	file.Pf("opts = append(")
	file.Pf("opts, kit_http.ServerBefore(")
	if helpers.IsTracingEnabled(service) && service.Options.Generate.Logger {
		file.Pf("opentracing.HTTPToContext(tracer, method, logger),")
	}
	file.Pf("goms_http.MethodInjector(\"%s\", method),", helpers.GetName(serviceName, service.Alias))
	file.Pf("goms_http.RequestIDCreator(),")
	file.Pf("goms_http.CorrelationIDExtractor(),")
	if helpers.IsLoggingEnabled(service) {
		file.Pf("goms_http.LoggerInjector(logger),")
	}
	file.Pf("),")
	file.Pf(")")
	file.Pf("return")
	file.Pf("},")
	file.Pf(")")
	file.Pf("")
	if service.Options.Generate.Logger {
		file.Pf("logger.Log(\"listening on\", addr)")
	}
	file.Pf("ch := make(chan error)")
	file.Pf("go func() {")
	file.Pf("ch <- server.ListenAndServe()")
	file.Pf("}()")
	file.Pf("select {")
	file.Pf("case err := <-ch:")
	file.Pf("if err == http.ErrServerClosed {")
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("return fmt.Errorf(\"http server: serve: %%v\", err)")
	file.Pf("case <-ctx.Done():")
	file.Pf("return server.Shutdown(context.Background())")
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}
