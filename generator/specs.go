package generator

import (
	"path/filepath"

	"github.com/wlMalk/goms/constants"
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/generators"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func DockerfileFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameDockerfile,
		file.NewSpec("Dockerfile").
			Name("Dockerfile", nil).
			Conditions(func(service types.Service) bool {
				return service.Generate.Has(constants.ServiceGenerateDockerfileFlag)
			}).
			Before(file.SpecBeforeFunc(func(file file.File, service types.Service) {
				file.(*files.TextFile).CommentFormat("# %s")
			})))
	g.AddServiceGenerator(constants.SpecNameDockerfile, constants.ServiceGeneratorDockerFileDefinition, generators.DockerFileDefinition)
}

func ProtoBufServiceDefinitionsFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameProtoBufServiceDefinitions,
		file.NewSpec("proto").
			Path("proto", nil).
			Name("service.goms", nil).
			Overwrite(true, nil).
			Conditions(func(service types.Service) bool {
				return service.Generate.Has(constants.ServiceGenerateProtoBufFlag) && helpers.IsGRPCEnabled(service)
			}).
			Before(file.SpecBeforeFunc(func(file file.File, service types.Service) {
				file.(*files.ProtoFile).Pkg = strings.ToLower(strings.ToSnakeCase(service.Name))
			})))
	g.AddServiceGenerator(constants.SpecNameProtoBufServiceDefinitions, constants.ServiceGeneratorProtoBufPackageDefinition, generators.ProtoBufPackageDefinition)
	g.AddMethodGeneratorWithExtractorAndConditions(constants.SpecNameProtoBufServiceDefinitions, constants.MethodGeneratorProtoBufMethodRequestDefinition, generators.ProtoBufMethodRequestDefinition, helpers.GetMethodsWithGRPCEnabled, func(service types.Service, method types.Method) bool {
		return len(method.Arguments) > 0
	})
	g.AddMethodGeneratorWithExtractorAndConditions(constants.SpecNameProtoBufServiceDefinitions, constants.MethodGeneratorProtoBufMethodResponseDefinition, generators.ProtoBufMethodResponseDefinition, helpers.GetMethodsWithGRPCEnabled, func(service types.Service, method types.Method) bool {
		return len(method.Results) > 0
	})
	g.AddServiceGenerator(constants.SpecNameProtoBufServiceDefinitions, constants.ServiceGeneratorProtoBufServiceDefinition, generators.ProtoBufServiceDefinition)
	g.AddEntityGenerator(constants.SpecNameProtoBufServiceDefinitions, constants.EntityGeneratorProtoBufEntityDefinition, generators.ProtoBufEntityDefinition)
	g.AddArgumentsGroupGenerator(constants.SpecNameProtoBufServiceDefinitions, constants.ArgumentsGroupGeneratorProtoBufArgumentsGroupDefinition, generators.ProtoBufArgumentsGroupDefinition)
	g.AddEnumGenerator(constants.SpecNameProtoBufServiceDefinitions, constants.EnumGeneratorProtoBufEnumDefinition, generators.ProtoBufEnumDefinition)
}

func ServiceMainFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceMain,
		file.NewSpec("go").
			Name("main", nil).
			Conditions(func(service types.Service) bool {
				return service.Generate.Has(constants.ServiceGenerateMainFlag)
			}).
			Before(file.SpecBeforeFunc(func(file file.File, service types.Service) {
				file.(*files.GoFile).Pkg = "main"
			})))
	g.AddServiceGenerator(constants.SpecNameServiceMain, constants.ServiceGeneratorServiceMainFunc, generators.ServiceMainFunc)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceMain, constants.ServiceGeneratorServiceMainInitLoggerFunc, generators.ServiceMainInitLoggerFunc, func(service types.Service) bool {
		return service.Generate.Has(constants.ServiceGenerateLoggerFlag) || helpers.IsLoggingEnabled(service)
	})
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceMain, constants.ServiceGeneratorServiceMainInitTracerFunc, generators.ServiceMainInitTracerFunc, helpers.IsTracingEnabled)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceMain, constants.ServiceGeneratorServiceMainInitFrequencyFunc, generators.ServiceMainInitFrequencyFunc, helpers.IsFrequencyMetricEnabled)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceMain, constants.ServiceGeneratorServiceMainInitLatencyFunc, generators.ServiceMainInitLatencyFunc, helpers.IsLatencyMetricEnabled)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceMain, constants.ServiceGeneratorServiceMainInitCounterFunc, generators.ServiceMainInitCounterFunc, helpers.IsCounterMetricEnabled)
}

func ServiceStartCMDFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceStartCMD,
		file.NewSpec("go").
			Path(filepath.Join("cmd", "start"), nil).
			Name("start", nil).
			Overwrite(true, nil).
			Conditions(func(service types.Service) bool {
				return service.Generate.Has(constants.ServiceGenerateMainFlag) && helpers.IsServerEnabled(service)
			}))
	g.AddServiceGenerator(constants.SpecNameServiceStartCMD, constants.ServiceGeneratorServiceStartCMDFunc, generators.ServiceStartCMDFunc)
	g.AddServiceGenerator(constants.SpecNameServiceStartCMD, constants.ServiceGeneratorServiceMainInitEndpointsFunc, generators.ServiceMainInitEndpointsFunc)
	g.AddServiceGenerator(constants.SpecNameServiceStartCMD, constants.ServiceGeneratorServiceMainPrepareEndpointsFunc, generators.ServiceMainPrepareEndpointsFunc)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceStartCMD, constants.ServiceGeneratorServiceMainInterruptHandlerFunc, generators.ServiceMainInterruptHandlerFunc, helpers.IsServerEnabled)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceStartCMD, constants.ServiceGeneratorServiceMainServeGRPCFunc, generators.ServiceMainServeGRPCFunc, helpers.IsGRPCServerEnabled)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceStartCMD, constants.ServiceGeneratorServiceMainServeHTTPFunc, generators.ServiceMainServeHTTPFunc, helpers.IsHTTPServerEnabled)
}

func CachingMiddlewareFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameCachingMiddleware,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "middleware"), nil).
			Name("caching_middleware.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsMiddlewareEnabled, helpers.IsCachingEnabled, helpers.IsCachaeble))
	g.AddServiceGenerator(constants.SpecNameCachingMiddleware, constants.ServiceGeneratorCachingMiddlewareStruct, generators.CachingMiddlewareStruct)
	g.AddServiceGenerator(constants.SpecNameCachingMiddleware, constants.ServiceGeneratorCachingMiddlewareCacheKeyerInterface, generators.CachingMiddlewareCacheKeyerInterface)
	g.AddServiceGenerator(constants.SpecNameCachingMiddleware, constants.ServiceGeneratorCachingMiddlewareNewFunc, generators.CachingMiddlewareNewFunc)
	g.AddMethodGenerator(constants.SpecNameCachingMiddleware, constants.MethodGeneratorCachingMiddlewareMethodFunc, generators.CachingMiddlewareMethodFunc)
}

func ConvertersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameConverters,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "handlers", "converters"), nil).
			Name("converters.goms", nil).
			Overwrite(true, nil))
	g.AddServiceGenerator(constants.SpecNameConverters, constants.ServiceGeneratorHandlerConverterTypes, generators.HandlerConverterTypes)
	g.AddServiceGenerator(constants.SpecNameConverters, constants.ServiceGeneratorHandlerConverterNewFuncs, generators.HandlerConverterNewFuncs)
	g.AddMethodGenerator(constants.SpecNameConverters, constants.MethodGeneratorHandlerToRequestResponseHandlerConverter, generators.HandlerToRequestResponseHandlerConverter)
	g.AddMethodGenerator(constants.SpecNameConverters, constants.MethodGeneratorRequestResponseHandlerToHandlerConverter, generators.RequestResponseHandlerToHandlerConverter)
	g.AddMethodGenerator(constants.SpecNameConverters, constants.MethodGeneratorRequestResponseHandlerToEndpointConverter, generators.RequestResponseHandlerToEndpointConverter)
	g.AddMethodGenerator(constants.SpecNameConverters, constants.MethodGeneratorEndpointToRequestResponseHandlerConverter, generators.EndpointToRequestResponseHandlerConverter)
}

func HandlersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHandlers,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "handlers"), nil).
			Name("handlers.goms", nil).
			Overwrite(true, nil))
	g.AddServiceGenerator(constants.SpecNameHandlers, constants.ServiceGeneratorServiceHandlerTypes, generators.ServiceHandlerTypes)
	g.AddMethodGenerator(constants.SpecNameHandlers, constants.MethodGeneratorMethodHandlers, generators.MethodHandlers)
}

func ServiceImplementationFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceImplementation,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return strings.ToLower(strings.ToSnakeCase(service.Name))
			}).
			Name("", func(service types.Service) string {
				return strings.ToLower(strings.ToSnakeCase(service.Name))
			}).
			Conditions(helpers.IsMethodStubsEnabled))
	g.AddServiceGenerator(constants.SpecNameServiceImplementation, constants.ServiceGeneratorServiceImplementationStruct, generators.ServiceImplementationStruct)
	g.AddServiceGenerator(constants.SpecNameServiceImplementation, constants.ServiceGeneratorServiceImplementationStructNewFunc, generators.ServiceImplementationStructNewFunc)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameServiceImplementation, constants.MethodGeneratorServiceMethodImplementation, generators.ServiceMethodImplementation, helpers.GetMethodsWithMethodStubsEnabled)
}

func ServiceImplementationValidatorFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceImplementationValidator,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return strings.ToLower(strings.ToSnakeCase(service.Name))
			}).
			Name("validator", nil).
			Conditions(helpers.IsMiddlewareEnabled, helpers.IsValidatingEnabled, helpers.IsValidatable))
	g.AddServiceGenerator(constants.SpecNameServiceImplementationValidator, constants.MethodGeneratorServiceMethodImplementationValidatorStruct, generators.ServiceMethodImplementationValidatorStruct)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameServiceImplementationValidator, constants.MethodGeneratorServiceMethodImplementationValidateFunc, generators.ServiceMethodImplementationValidateFunc, helpers.GetMethodsWithValidatingEnabled)
}

func ServiceImplementationMiddlewareFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceImplementationMiddleware,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return strings.ToLower(strings.ToSnakeCase(service.Name))
			}).
			Name("middleware", nil).
			Conditions(helpers.IsMiddlewareEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameServiceImplementationMiddleware, constants.MethodGeneratorServiceMethodImplementationMiddleware, generators.ServiceMethodImplementationMiddleware, helpers.GetMethodsWithMiddlewareEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameServiceImplementationMiddleware, constants.MethodGeneratorServiceMethodImplementationOuterMiddleware, generators.ServiceMethodImplementationOuterMiddleware, helpers.GetMethodsWithMiddlewareEnabled)
	g.AddServiceGeneratorWithConditions(constants.SpecNameServiceImplementationMiddleware, constants.ServiceGeneratorServiceImplementationMiddleware, generators.ServiceImplementationMiddleware, func(service types.Service) bool {
		return service.Generate.Has(constants.ServiceGenerateMiddlewareFlag)
	})
}

func CachingKeyerFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameCachingKeyer,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return strings.ToLower(strings.ToSnakeCase(service.Name))
			}).
			Name("caching_keyer", nil).
			Conditions(helpers.IsMiddlewareEnabled, helpers.IsCachingEnabled, helpers.IsCachaeble))
	g.AddServiceGenerator(constants.SpecNameCachingKeyer, constants.ServiceGeneratorCachingMiddlewareCacheKeyerType, generators.CachingMiddlewareCacheKeyerType)
	g.AddServiceGenerator(constants.SpecNameCachingKeyer, constants.ServiceGeneratorCachingMiddlewareKeyerNewFunc, generators.CachingMiddlewareKeyerNewFunc)
	g.AddMethodGeneratorWithExtractorAndConditions(constants.SpecNameCachingKeyer, constants.MethodGeneratorCachingMiddlewareKeyerMethodFunc, generators.CachingMiddlewareKeyerMethodFunc, helpers.GetMethodsWithCachingEnabled, func(service types.Service, method types.Method) bool {
		return len(method.Arguments) > 0 && len(method.Results) > 0
	})
}

func LoggingMiddlewareFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameLoggingMiddleware,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "middleware"), nil).
			Name("logging_middleware.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsMiddlewareEnabled, helpers.IsLoggingEnabled))
	g.AddServiceGeneratorWithConditions(constants.SpecNameLoggingMiddleware, constants.ServiceGeneratorLoggingMiddlewareStructs, generators.LoggingMiddlewareStructs, func(service types.Service) bool {
		return helpers.HasLoggeds(service) || helpers.HasLoggedErrors(service)
	})
	g.AddServiceGeneratorWithConditions(constants.SpecNameLoggingMiddleware, constants.ServiceGeneratorLoggingMiddlewareNewFunc, generators.LoggingMiddlewareNewFunc, func(service types.Service) bool {
		return helpers.HasLoggeds(service) || helpers.HasLoggedErrors(service)
	})
	g.AddMethodGeneratorWithConditions(constants.SpecNameLoggingMiddleware, constants.MethodGeneratorLoggingMiddlewareMethodHandler, generators.LoggingMiddlewareMethodHandler, func(service types.Service, _ types.Method) bool {
		return helpers.HasLoggeds(service)
	})
	g.AddMethodGeneratorWithConditions(constants.SpecNameLoggingMiddleware, constants.MethodGeneratorErrorLoggingMiddlewareMethodHandler, generators.ErrorLoggingMiddlewareMethodHandler, func(service types.Service, _ types.Method) bool {
		return helpers.HasLoggedErrors(service)
	})
	g.AddServiceGeneratorWithConditions(constants.SpecNameLoggingMiddleware, constants.ServiceGeneratorLoggingMiddlewareTypes, generators.LoggingMiddlewareTypes, helpers.HasLoggeds)
}

func ServiceMiddlewareFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceMiddleware,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "middleware"), nil).
			Name("middleware.goms", nil).
			Overwrite(true, nil).
			Conditions(func(service types.Service) bool {
				return helpers.IsMiddlewareEnabled(service) || service.Generate.Has(constants.ServiceGenerateMiddlewareFlag)
			}))
	g.AddServiceGenerator(constants.SpecNameServiceMiddleware, constants.ServiceGeneratorServiceMiddlewareTypes, generators.ServiceMiddlewareTypes)
	g.AddServiceGenerator(constants.SpecNameServiceMiddleware, constants.ServiceGeneratorServiceMiddlewareChainFunc, generators.ServiceMiddlewareChainFunc)
	g.AddServiceGenerator(constants.SpecNameServiceMiddleware, constants.ServiceGeneratorServiceRequestResponseMiddlewareChainFunc, generators.ServiceRequestResponseMiddlewareChainFunc)
	g.AddServiceGenerator(constants.SpecNameServiceMiddleware, constants.ServiceGeneratorServiceApplyMiddlewareFunc, generators.ServiceApplyMiddlewareFunc)
	g.AddServiceGenerator(constants.SpecNameServiceMiddleware, constants.ServiceGeneratorServiceApplyMiddlewareSpecialFunc, generators.ServiceApplyMiddlewareSpecialFunc)
	g.AddServiceGenerator(constants.SpecNameServiceMiddleware, constants.ServiceGeneratorServiceApplyMiddlewareConditionalFunc, generators.ServiceApplyMiddlewareConditionalFunc)
}

func RecoveringMiddlewareFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameRecoveringMiddleware,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "middleware"), nil).
			Name("recovering_middleware.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsMiddlewareEnabled, helpers.IsRecoveringEnabled))
	g.AddServiceGenerator(constants.SpecNameRecoveringMiddleware, constants.ServiceGeneratorRecoveringMiddlewareStruct, generators.RecoveringMiddlewareStruct)
	g.AddServiceGenerator(constants.SpecNameRecoveringMiddleware, constants.ServiceGeneratorRecoveringMiddlewareNewFunc, generators.RecoveringMiddlewareNewFunc)
	g.AddMethodGenerator(constants.SpecNameRecoveringMiddleware, constants.MethodGeneratorRecoveringMiddlewareMethodFunc, generators.RecoveringMiddlewareMethodFunc)
}

func ServiceTransportEndpointsFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameServiceTransportEndpoints,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport"), nil).
			Name("endpoints.goms", nil).
			Overwrite(true, nil))
	g.AddServiceGenerator(constants.SpecNameServiceTransportEndpoints, constants.ServiceGeneratorServiceStructType, generators.ServiceStructType)
	g.AddServiceGenerator(constants.SpecNameServiceTransportEndpoints, constants.ServiceGeneratorServiceStructTypeNewFunc, generators.ServiceStructTypeNewFunc)
	g.AddMethodGenerator(constants.SpecNameServiceTransportEndpoints, constants.MethodGeneratorServiceStructMethodHandler, generators.ServiceStructMethodHandler)
}

func ValidatingMiddlewareFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameValidatingMiddleware,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "middleware"), nil).
			Name("validating_middleware.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsMiddlewareEnabled, helpers.IsValidatingEnabled, helpers.IsValidatable))
	g.AddServiceGenerator(constants.SpecNameValidatingMiddleware, constants.ServiceGeneratorValidatingValidatorsTypes, generators.ValidatingValidatorsTypes)
	g.AddServiceGenerator(constants.SpecNameValidatingMiddleware, constants.ServiceGeneratorValidatingMiddlewareStruct, generators.ValidatingMiddlewareStruct)
	g.AddServiceGenerator(constants.SpecNameValidatingMiddleware, constants.ServiceGeneratorValidatingMiddlewareNewFunc, generators.ValidatingMiddlewareNewFunc)
	g.AddMethodGenerator(constants.SpecNameValidatingMiddleware, constants.MethodGeneratorValidatingMiddlewareMethodFunc, generators.ValidatingMiddlewareMethodFunc)
}

func GRPCClientFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameGRPCClient,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "grpc", "client"), nil).
			Name("client.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsGRPCClientEnabled))
	g.AddServiceGenerator(constants.SpecNameGRPCClient, constants.ServiceGeneratorGRPCTransportClientStruct, generators.GRPCTransportClientStruct)
	g.AddServiceGenerator(constants.SpecNameGRPCClient, constants.ServiceGeneratorGRPCTransportClientNewFunc, generators.GRPCTransportClientNewFunc)
	g.AddServiceGenerator(constants.SpecNameGRPCClient, constants.ServiceGeneratorGRPCTransportClientNewSpecialFunc, generators.GRPCTransportClientNewSpecialFunc)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGRPCClient, constants.MethodGeneratorGRPCTransportClientMethodFunc, generators.GRPCTransportClientMethodFunc, helpers.GetMethodsWithGRPCClientEnabled)
}

func GlobalGRPCClientFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameGlobalGRPCClient,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return filepath.Join("clients", "grpc", strings.ToLower(strings.ToSnakeCase(service.Name)))
			}).
			Name("client.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsGRPCClientEnabled))
	g.AddServiceGenerator(constants.SpecNameGlobalGRPCClient, constants.ServiceGeneratorGRPCTransportClientGlobalVar, generators.GRPCTransportClientGlobalVar)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGlobalGRPCClient, constants.MethodGeneratorGRPCTransportClientGlobalFunc, generators.GRPCTransportClientGlobalFunc, helpers.GetMethodsWithGRPCClientEnabled)
}

func GRPCDecodersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameGRPCDecoders,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "grpc"), nil).
			Name("decoders.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsGRPCEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGRPCDecoders, constants.MethodGeneratorGRPCRequestDecoder, generators.GRPCRequestDecoder, helpers.GetMethodsWithGRPCEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGRPCDecoders, constants.MethodGeneratorGRPCResponseDecoder, generators.GRPCResponseDecoder, helpers.GetMethodsWithGRPCEnabled)
}

func GRPCEncodersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameGRPCEncoders,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "grpc"), nil).
			Name("encoders.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsGRPCEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGRPCEncoders, constants.MethodGeneratorGRPCRequestEncoder, generators.GRPCRequestEncoder, helpers.GetMethodsWithGRPCEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGRPCEncoders, constants.MethodGeneratorGRPCResponseEncoder, generators.GRPCResponseEncoder, helpers.GetMethodsWithGRPCEnabled)
}

func GRPCServerFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameGRPCServer,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "grpc", "server"), nil).
			Name("server.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsGRPCServerEnabled))
	g.AddServiceGenerator(constants.SpecNameGRPCServer, constants.ServiceGeneratorGRPCTransportServerRegisterFunc, generators.GRPCTransportServerRegisterFunc)
	g.AddServiceGenerator(constants.SpecNameGRPCServer, constants.ServiceGeneratorGRPCTransportServerRegisterSpecialFunc, generators.GRPCTransportServerRegisterSpecialFunc)
	g.AddServiceGenerator(constants.SpecNameGRPCServer, constants.ServiceGeneratorGRPCTransportServerHandlerStruct, generators.GRPCTransportServerHandlerStruct)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGRPCServer, constants.MethodGeneratorGRPCTransportServerHandlerMethodFunc, generators.GRPCTransportServerHandlerMethodFunc, helpers.GetMethodsWithGRPCServerEnabled)
}

func HTTPClientFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHTTPClient,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "http", "client"), nil).
			Name("client.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPClientEnabled))
	g.AddServiceGenerator(constants.SpecNameHTTPClient, constants.ServiceGeneratorHTTPTransportClientStruct, generators.HTTPTransportClientStruct)
	g.AddServiceGenerator(constants.SpecNameHTTPClient, constants.ServiceGeneratorHTTPTransportClientNewFunc, generators.HTTPTransportClientNewFunc)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPClient, constants.MethodGeneratorHTTPTransportClientMethodFunc, generators.HTTPTransportClientMethodFunc, helpers.GetMethodsWithHTTPClientEnabled)
}

func GlobalHTTPClientFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameGlobalHTTPClient,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return filepath.Join("clients", "http", strings.ToLower(strings.ToSnakeCase(service.Name)))
			}).
			Name("client.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPClientEnabled))
	g.AddServiceGenerator(constants.SpecNameGlobalHTTPClient, constants.ServiceGeneratorHTTPTransportClientGlobalVar, generators.HTTPTransportClientGlobalVar)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameGlobalHTTPClient, constants.MethodGeneratorHTTPTransportClientGlobalFunc, generators.HTTPTransportClientGlobalFunc, helpers.GetMethodsWithHTTPClientEnabled)
}

func HTTPDecodersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHTTPDecoders,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "http"), nil).
			Name("decoders.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPDecoders, constants.MethodGeneratorHTTPRequestDecoder, generators.HTTPRequestDecoder, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPDecoders, constants.MethodGeneratorHTTPResponseDecoder, generators.HTTPResponseDecoder, helpers.GetMethodsWithHTTPEnabled)
}

func HTTPEncodersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHTTPEncoders,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "http"), nil).
			Name("encoders.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPEncoders, constants.MethodGeneratorHTTPRequestEncoder, generators.HTTPRequestEncoder, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPEncoders, constants.MethodGeneratorHTTPResponseEncoder, generators.HTTPResponseEncoder, helpers.GetMethodsWithHTTPEnabled)
}

func HTTPRequestsFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHTTPRequests,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "http", "requests"), nil).
			Name("requests.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPRequests, constants.MethodGeneratorHTTPRequest, generators.HTTPRequest, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPRequests, constants.MethodGeneratorHTTPRequestNewFunc, generators.HTTPRequestNewFunc, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPRequests, constants.MethodGeneratorHTTPRequestNewHTTPFunc, generators.HTTPRequestNewHTTPFunc, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPRequests, constants.MethodGeneratorHTTPRequestToRequestFunc, generators.HTTPRequestToRequestFunc, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPRequests, constants.MethodGeneratorHTTPRequestToHTTPArgFunc, generators.HTTPRequestToHTTPArgFunc, helpers.GetMethodsWithHTTPEnabled)
}

func HTTPResponsesFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHTTPResponses,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "http", "responses"), nil).
			Name("responses.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPEnabled))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPResponses, constants.MethodGeneratorHTTPResponse, generators.HTTPResponse, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPResponses, constants.MethodGeneratorHTTPResponseNewFunc, generators.HTTPResponseNewFunc, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPResponses, constants.MethodGeneratorHTTPResponseNewHTTPFunc, generators.HTTPResponseNewHTTPFunc, helpers.GetMethodsWithHTTPEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameHTTPResponses, constants.MethodGeneratorHTTPResponseToResponseFunc, generators.HTTPResponseToResponseFunc, helpers.GetMethodsWithHTTPEnabled)
}

func HTTPServerFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameHTTPServer,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "transport", "http", "server"), nil).
			Name("server.goms", nil).
			Overwrite(true, nil).
			Conditions(helpers.IsHTTPServerEnabled))
	g.AddServiceGenerator(constants.SpecNameHTTPServer, constants.ServiceGeneratorHTTPTransportServerRegisterFunc, generators.HTTPTransportServerRegisterFunc)
	g.AddServiceGenerator(constants.SpecNameHTTPServer, constants.ServiceGeneratorHTTPTransportServerRegisterSpecialFunc, generators.HTTPTransportServerRegisterSpecialFunc)
}

func ProtoRequestsConvertersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameProtoRequestsConverters,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return filepath.Join("pkg", "protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)), "requests")
			}).
			Name("requests.goms", nil).
			Overwrite(true, nil).
			Conditions(func(service types.Service) bool {
				return service.Generate.Has(constants.ServiceGenerateProtoBufFlag) && helpers.IsGRPCEnabled(service)
			}))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameProtoRequestsConverters, constants.MethodGeneratorProtoRequestNewFunc, generators.ProtoRequestNewFunc, helpers.GetMethodsWithGRPCEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameProtoRequestsConverters, constants.MethodGeneratorProtoRequestNewProtoFunc, generators.ProtoRequestNewProtoFunc, helpers.GetMethodsWithGRPCEnabled)
}

func ProtoResponsesConvertersFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameProtoResponsesConverters,
		file.NewSpec("go").
			Path("", func(service types.Service) string {
				return filepath.Join("pkg", "protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)), "responses")
			}).
			Name("responses.goms", nil).
			Overwrite(true, nil).
			Conditions(func(service types.Service) bool {
				return service.Generate.Has(constants.ServiceGenerateProtoBufFlag) && helpers.IsGRPCEnabled(service)
			}))
	g.AddMethodGeneratorWithExtractor(constants.SpecNameProtoResponsesConverters, constants.MethodGeneratorProtoResponseNewFunc, generators.ProtoResponseNewFunc, helpers.GetMethodsWithGRPCEnabled)
	g.AddMethodGeneratorWithExtractor(constants.SpecNameProtoResponsesConverters, constants.MethodGeneratorProtoResponseNewProtoFunc, generators.ProtoResponseNewProtoFunc, helpers.GetMethodsWithGRPCEnabled)
}

func RequestsFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameRequests,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "requests"), nil).
			Name("requests.goms", nil).
			Overwrite(true, nil))
	g.AddMethodGeneratorWithConditions(constants.SpecNameRequests, constants.MethodGeneratorServiceRequestStruct, generators.ServiceRequestStruct, func(service types.Service, method types.Method) bool {
		return len(method.Arguments) > 0
	})
	g.AddMethodGeneratorWithConditions(constants.SpecNameRequests, constants.MethodGeneratorServiceRequestNewFunc, generators.ServiceRequestNewFunc, func(service types.Service, method types.Method) bool {
		return len(method.Arguments) > 0
	})
}

func ResponseFileSpec(g *Generator) {
	g.AddSpec(constants.SpecNameResponses,
		file.NewSpec("go").
			Path(filepath.Join("pkg", "service", "responses"), nil).
			Name("responses.goms", nil).
			Overwrite(true, nil))
	g.AddMethodGeneratorWithConditions(constants.SpecNameResponses, constants.MethodGeneratorServiceResponseStruct, generators.ServiceResponseStruct, func(service types.Service, method types.Method) bool {
		return len(method.Results) > 0
	})
}
