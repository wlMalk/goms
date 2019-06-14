package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateServiceStartCMDFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateServiceStartCMDFunc(file, service)
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

func generateServiceStartCMDFunc(file *files.GoFile, service *types.Service) {
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("", "context")
	file.AddImport("", "errors")
	file.AddImport("", "fmt")
	file.AddImport("", "net/http")
	file.AddImport("", "os")
	file.AddImport("", "os/signal")
	file.AddImport("", "syscall")
	file.AddImport("", "time")

	file.AddImport("", "github.com/go-kit/kit/endpoint")

	file.AddImport("", service.ImportPath, "/pkg/service/handlers/converters")
	file.AddImport("", service.ImportPath, "/pkg/service/transport")
	file.AddImport("", service.ImportPath, "/"+strings.ToSnakeCase(service.Name))

	if helpers.IsServerEnabled(service) {
		file.AddImport("", "golang.org/x/sync/errgroup")
	}
	if helpers.IsRateLimitingEnabled(service) ||
		helpers.IsCircuitBreakingEnabled(service) ||
		helpers.IsRecoveringEnabled(service) ||
		helpers.IsLoggingEnabled(service) ||
		helpers.IsTracingEnabled(service) ||
		helpers.IsMetricsEnabled(service) {
		file.AddImport("", service.ImportPath, "/pkg/service/middleware")
	}
	if helpers.IsRecoveringEnabled(service) ||
		helpers.IsLoggingEnabled(service) ||
		helpers.IsMetricsEnabled(service) {
		file.AddImport("goms_middleware", "github.com/wlMalk/goms/goms/middleware")
	}
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		file.AddImport("", "github.com/go-kit/kit/log")
	}
	if helpers.IsMetricsEnabled(service) {
		file.AddImport("", "github.com/go-kit/kit/metrics")
	}
	if helpers.IsTracingEnabled(service) {
		file.AddImport("", "github.com/go-kit/kit/tracing/opentracing")
		file.AddImport("opentracinggo", "github.com/opentracing/opentracing-go")
	}
	if helpers.IsRateLimitingEnabled(service) {
		file.AddImport("", "golang.org/x/time/rate")
		file.AddImport("", "github.com/go-kit/kit/ratelimit")
	}
	if helpers.IsCircuitBreakingEnabled(service) {
		file.AddImport("", "github.com/go-kit/kit/circuitbreaker")
		file.AddImport("", "github.com/sony/gobreaker")
	}
	if helpers.IsHTTPServerEnabled(service) {
		file.AddImport(strings.ToSnakeCase(service.Name)+"_http_server", service.ImportPath, "/pkg/service/transport/http/server")
		file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
		file.AddImport("goms_router", "github.com/wlMalk/goms/goms/transport/http/httprouter")
		file.AddImport("", "github.com/julienschmidt/httprouter")
	}

	file.Pf("func Start(")
	if service.Options.Generate.Logger || helpers.IsLoggingEnabled(service) {
		file.Pf("logger log.Logger,")
	}
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
	file.Pf(") {")
	if service.Options.Generate.Logger {
		file.Pf("logger.Log(\"message\", \"Hello, I am alive\")")
		file.Pf("defer logger.Log(\"message\", \"goodbye, good luck\")")
		file.Pf("")
	}

	if helpers.IsServerEnabled(service) {
		file.Pf("g, ctx := errgroup.WithContext(context.Background())")
		file.Pf("g.Go(func() error {")
		file.Pf("return interruptHandler(ctx)")
		file.Pf("})")
		file.Pf("")
	}
	file.Pf("s := %s.New()", serviceNameSnake)
	file.Pf("endpoints := initEndpoints(s)")
	file.Pf("endpoints = prepareEndpoints(")
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
		file.Pf("return serveHTTP(")
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

func generateServiceMainInitEndpointsFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.Pf("func initEndpoints(s *%s.%s) transport.%s {", serviceNameSnake, serviceName, serviceName)
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
	file.Pf("func prepareEndpoints(")
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
	file.Pf("func interruptHandler(ctx context.Context) error {")
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
	file.Pf("func serveHTTP(")
	file.Pf("ctx context.Context,")
	file.Pf("endpoints *transport.%s,", serviceName)
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
