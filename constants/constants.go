package constants

const (
	ServiceGeneratorCachingMiddlewareCacheKeyerInterface      string = "caching-middleware-cache-keyer-interface"
	ServiceGeneratorCachingMiddlewareCacheKeyerType           string = "caching-middleware-cache-keyer-type"
	ServiceGeneratorCachingMiddlewareKeyerNewFunc             string = "caching-middleware-keyer-new-func"
	ServiceGeneratorCachingMiddlewareNewFunc                  string = "caching-middleware-new-func"
	ServiceGeneratorCachingMiddlewareStruct                   string = "caching-middleware-struct"
	ServiceGeneratorDockerFileDefinition                      string = "docker-file-definition"
	ServiceGeneratorGRPCTransportClientGlobalVar              string = "grpc-transport-client-global-var"
	ServiceGeneratorGRPCTransportClientNewFunc                string = "grpc-transport-client-new-func"
	ServiceGeneratorGRPCTransportClientNewSpecialFunc         string = "grpc-transport-client-new-special-func"
	ServiceGeneratorGRPCTransportClientStruct                 string = "grpc-transport-client-struct"
	ServiceGeneratorGRPCTransportServerHandlerStruct          string = "grpc-transport-server-handler-struct"
	ServiceGeneratorGRPCTransportServerRegisterFunc           string = "grpc-transport-server-register-func"
	ServiceGeneratorGRPCTransportServerRegisterSpecialFunc    string = "grpc-transport-server-register-special-func"
	ServiceGeneratorHTTPTransportClientGlobalVar              string = "http-transport-client-global-var"
	ServiceGeneratorHTTPTransportClientNewFunc                string = "http-transport-client-new-func"
	ServiceGeneratorHTTPTransportClientNewSpecialFunc         string = "http-transport-client-new-special-func"
	ServiceGeneratorHTTPTransportClientStruct                 string = "http-transport-client-struct"
	ServiceGeneratorHTTPTransportServerRegisterFunc           string = "http-transport-server-register-func"
	ServiceGeneratorHTTPTransportServerRegisterSpecialFunc    string = "http-transport-server-register-special-func"
	ServiceGeneratorHandlerConverterNewFuncs                  string = "handler-converter-new-funcs"
	ServiceGeneratorHandlerConverterTypes                     string = "handler-converter-types"
	ServiceGeneratorLocalClientGlobalVar                      string = "local-client-global-var"
	ServiceGeneratorLoggingMiddlewareNewFunc                  string = "logging-middleware-new-func"
	ServiceGeneratorLoggingMiddlewareStructs                  string = "logging-middleware-structs"
	ServiceGeneratorLoggingMiddlewareTypes                    string = "logging-middleware-types"
	ServiceGeneratorProtoBufPackageDefinition                 string = "proto-buf-package-definition"
	ServiceGeneratorProtoBufServiceDefinition                 string = "proto-buf-service-definition"
	ServiceGeneratorRecoveringMiddlewareNewFunc               string = "recovering-middleware-new-func"
	ServiceGeneratorRecoveringMiddlewareStruct                string = "recovering-middleware-struct"
	ServiceGeneratorServiceApplyMiddlewareConditionalFunc     string = "service-apply-middleware-conditional-func"
	ServiceGeneratorServiceApplyMiddlewareFunc                string = "service-apply-middleware-func"
	ServiceGeneratorServiceApplyMiddlewareSpecialFunc         string = "service-apply-middleware-special-func"
	ServiceGeneratorServiceHandlerTypes                       string = "service-handler-types"
	ServiceGeneratorServiceImplementationMiddleware           string = "service-implementation-middleware"
	ServiceGeneratorServiceImplementationStruct               string = "service-implementation-struct"
	ServiceGeneratorServiceImplementationStructNewFunc        string = "service-implementation-struct-new-func"
	ServiceGeneratorServiceMainFunc                           string = "service-main-func"
	ServiceGeneratorServiceMainInitCounterFunc                string = "service-main-init-counter-func"
	ServiceGeneratorServiceMainInitEndpointsFunc              string = "service-main-init-endpoints-func"
	ServiceGeneratorServiceMainInitFrequencyFunc              string = "service-main-init-frequency-func"
	ServiceGeneratorServiceMainInitLatencyFunc                string = "service-main-init-latency-func"
	ServiceGeneratorServiceMainInitLoggerFunc                 string = "service-main-init-logger-func"
	ServiceGeneratorServiceMainInitTracerFunc                 string = "service-main-init-tracer-func"
	ServiceGeneratorServiceMainInterruptHandlerFunc           string = "service-main-interrupt-handler-func"
	ServiceGeneratorServiceMainPrepareEndpointsFunc           string = "service-main-prepare-endpoints-func"
	ServiceGeneratorServiceMainServeGRPCFunc                  string = "service-main-serve-grpc-func"
	ServiceGeneratorServiceMainServeHTTPFunc                  string = "service-main-serve-http-func"
	ServiceGeneratorServiceMiddlewareChainFunc                string = "service-middleware-chain-func"
	ServiceGeneratorServiceMiddlewareTypes                    string = "service-middleware-types"
	ServiceGeneratorServiceRequestResponseMiddlewareChainFunc string = "service-request-response-middleware-chain-func"
	ServiceGeneratorServiceStartCMDFunc                       string = "service-start-cmd-func"
	ServiceGeneratorServiceStructType                         string = "service-struct-type"
	ServiceGeneratorServiceStructTypeNewFunc                  string = "service-struct-type-new-func"
	ServiceGeneratorValidatingMiddlewareNewFunc               string = "validating-middleware-new-func"
	ServiceGeneratorValidatingMiddlewareStruct                string = "validating-middleware-struct"
	ServiceGeneratorValidatingValidatorsTypes                 string = "validating-validators-types"
)

const (
	MethodGeneratorCachingMiddlewareKeyerMethodFunc           string = "caching-middleware-keyer-method-func"
	MethodGeneratorCachingMiddlewareMethodFunc                string = "caching-middleware-method-func"
	MethodGeneratorEndpointToRequestResponseHandlerConverter  string = "endpoint-to-request-response-handler-converter"
	MethodGeneratorErrorLoggingMiddlewareMethodHandler        string = "error-logging-middleware-method-handler"
	MethodGeneratorGRPCRequestDecoder                         string = "grpc-request-decoder"
	MethodGeneratorGRPCRequestEncoder                         string = "grpc-request-encoder"
	MethodGeneratorGRPCResponseDecoder                        string = "grpc-response-decoder"
	MethodGeneratorGRPCResponseEncoder                        string = "grpc-response-encoder"
	MethodGeneratorGRPCTransportClientGlobalFunc              string = "grpc-transport-client-global-func"
	MethodGeneratorGRPCTransportClientMethodFunc              string = "grpc-transport-client-method-func"
	MethodGeneratorGRPCTransportServerHandlerMethodFunc       string = "grpc-transport-server-handler-method-func"
	MethodGeneratorHTTPRequest                                string = "http-request"
	MethodGeneratorHTTPRequestDecoder                         string = "http-request-decoder"
	MethodGeneratorHTTPRequestEncoder                         string = "http-request-encoder"
	MethodGeneratorHTTPRequestNewFunc                         string = "http-request-new-func"
	MethodGeneratorHTTPRequestNewHTTPFunc                     string = "http-request-new-http-func"
	MethodGeneratorHTTPRequestToHTTPArgFunc                   string = "http-request-to-http-arg-func"
	MethodGeneratorHTTPRequestToRequestFunc                   string = "http-request-to-request-func"
	MethodGeneratorHTTPResponse                               string = "http-response"
	MethodGeneratorHTTPResponseDecoder                        string = "http-response-decoder"
	MethodGeneratorHTTPResponseEncoder                        string = "http-response-encoder"
	MethodGeneratorHTTPResponseNewFunc                        string = "http-response-new-func"
	MethodGeneratorHTTPResponseNewHTTPFunc                    string = "http-response-new-http-func"
	MethodGeneratorHTTPResponseToResponseFunc                 string = "http-response-to-response-func"
	MethodGeneratorHTTPTransportClientGlobalFunc              string = "http-transport-client-global-func"
	MethodGeneratorHTTPTransportClientMethodFunc              string = "http-transport-client-method-func"
	MethodGeneratorHandlerToRequestResponseHandlerConverter   string = "handler-to-request-response-handler-converter"
	MethodGeneratorLocalClientGlobalFunc                      string = "local-client-global-func"
	MethodGeneratorLoggingMiddlewareMethodHandler             string = "logging-middleware-method-handler"
	MethodGeneratorMethodHandlers                             string = "method-handlers"
	MethodGeneratorProtoBufMethodRequestDefinition            string = "proto-buf-method-request-definition"
	MethodGeneratorProtoBufMethodResponseDefinition           string = "proto-buf-method-response-definition"
	MethodGeneratorProtoRequestNewFunc                        string = "proto-request-new-func"
	MethodGeneratorProtoRequestNewProtoFunc                   string = "proto-request-new-proto-func"
	MethodGeneratorProtoResponseNewFunc                       string = "proto-response-new-func"
	MethodGeneratorProtoResponseNewProtoFunc                  string = "proto-response-new-proto-func"
	MethodGeneratorRecoveringMiddlewareMethodFunc             string = "recovering-middleware-method-func"
	MethodGeneratorRequestResponseHandlerToEndpointConverter  string = "request-response-handler-to-endpoint-converter"
	MethodGeneratorRequestResponseHandlerToHandlerConverter   string = "request-response-handler-to-handler-converter"
	MethodGeneratorServiceMethodImplementation                string = "service-method-implementation"
	MethodGeneratorServiceMethodImplementationMiddleware      string = "service-method-implementation-middleware"
	MethodGeneratorServiceMethodImplementationOuterMiddleware string = "service-method-implementation-outer-middleware"
	MethodGeneratorServiceMethodImplementationValidateFunc    string = "service-method-implementation-validate-func"
	MethodGeneratorServiceMethodImplementationValidatorStruct string = "service-method-implementation-validator-struct"
	MethodGeneratorServiceRequestNewFunc                      string = "service-request-new-func"
	MethodGeneratorServiceRequestStruct                       string = "service-request-struct"
	MethodGeneratorServiceResponseStruct                      string = "service-response-struct"
	MethodGeneratorServiceStructMethodHandler                 string = "service-struct-method-handler"
	MethodGeneratorValidatingMiddlewareMethodFunc             string = "validating-middleware-method-func"
)

const (
	EntityGeneratorProtoBufEntityDefinition string = "proto-buf-entity-definition"
)

const (
	ArgumentsGroupGeneratorProtoBufArgumentsGroupDefinition string = "proto-buf-arguments-group-definition"
)

const (
	EnumGeneratorProtoBufEnumDefinition string = "proto-buf-enum-definition"
)

const (
	SpecNameCachingKeyer                    string = "caching-keyer"
	SpecNameCachingMiddleware               string = "caching-middleware"
	SpecNameConverters                      string = "converters"
	SpecNameDockerfile                      string = "dockerfile"
	SpecNameGRPCClient                      string = "grpc-client"
	SpecNameGRPCDecoders                    string = "grpc-decoders"
	SpecNameGRPCEncoders                    string = "grpc-encoders"
	SpecNameGRPCServer                      string = "grpc-server"
	SpecNameGlobalGRPCClient                string = "global-grpc-client"
	SpecNameGlobalHTTPClient                string = "global-http-client"
	SpecNameGlobalLocalClient               string = "global-local-client"
	SpecNameHTTPClient                      string = "http-client"
	SpecNameHTTPDecoders                    string = "http-decoders"
	SpecNameHTTPEncoders                    string = "http-encoders"
	SpecNameHTTPRequests                    string = "http-requests"
	SpecNameHTTPResponses                   string = "http-responses"
	SpecNameHTTPServer                      string = "http-server"
	SpecNameHandlers                        string = "handlers"
	SpecNameLoggingMiddleware               string = "logging-middleware"
	SpecNameProtoBufServiceDefinitions      string = "proto-buf-service-definitions"
	SpecNameProtoRequestsConverters         string = "proto-requests-converters"
	SpecNameProtoResponsesConverters        string = "proto-responses-converters"
	SpecNameRecoveringMiddleware            string = "recovering-middleware"
	SpecNameRequests                        string = "requests"
	SpecNameResponses                       string = "responses"
	SpecNameServiceImplementation           string = "service-implementation"
	SpecNameServiceImplementationMiddleware string = "service-implementation-middleware"
	SpecNameServiceImplementationValidator  string = "service-implementation-validator"
	SpecNameServiceMain                     string = "service-main"
	SpecNameServiceMiddleware               string = "service-middleware"
	SpecNameServiceStartCMD                 string = "service-start-cmd"
	SpecNameServiceTransportEndpoints       string = "service-transport-endpoints"
	SpecNameValidatingMiddleware            string = "validating-middleware"
)

const (
	ServiceGenerateCachingFlag          string = "caching"
	ServiceGenerateCircuitBreakingFlag  string = "circuit-breaking"
	ServiceGenerateCounterMetricFlag    string = "counter-metric"
	ServiceGenerateDockerfileFlag       string = "dockerfile"
	ServiceGenerateFrequencyMetricFlag  string = "frequency-metric"
	ServiceGenerateGRPCClientFlag       string = "grpc-client"
	ServiceGenerateGRPCServerFlag       string = "grpc-server"
	ServiceGenerateHTTPClientFlag       string = "http-client"
	ServiceGenerateHTTPServerFlag       string = "http-server"
	ServiceGenerateLatencyMetricFlag    string = "latency-metric"
	ServiceGenerateLoggerFlag           string = "logger"
	ServiceGenerateLoggingFlag          string = "logging"
	ServiceGenerateMainFlag             string = "main"
	ServiceGenerateMethodStubsFlag      string = "method-stubs"
	ServiceGenerateMiddlewareFlag       string = "middleware"
	ServiceGenerateProtoBufFlag         string = "proto-buf"
	ServiceGenerateRateLimitingFlag     string = "rate-limiting"
	ServiceGenerateRecoveringFlag       string = "recovering"
	ServiceGenerateServiceDiscoveryFlag string = "service-discovery"
	ServiceGenerateTracingFlag          string = "tracing"
	ServiceGenerateValidatingFlag       string = "validating"
	ServiceGenerateValidatorsFlag       string = "validators"
)

const (
	MethodGenerateCachingFlag         string = "caching"
	MethodGenerateCircuitBreakingFlag string = "circuit-breaking"
	MethodGenerateCounterMetricFlag   string = "counter-metric"
	MethodGenerateFrequencyMetricFlag string = "frequency-metric"
	MethodGenerateGRPCClientFlag      string = "grpc-client"
	MethodGenerateGRPCServerFlag      string = "grpc-server"
	MethodGenerateHTTPClientFlag      string = "http-client"
	MethodGenerateHTTPServerFlag      string = "http-server"
	MethodGenerateLatencyMetricFlag   string = "latency-metric"
	MethodGenerateLoggingFlag         string = "logging"
	MethodGenerateMethodStubsFlag     string = "method-stubs"
	MethodGenerateMiddlewareFlag      string = "middleware"
	MethodGenerateRateLimitingFlag    string = "rate-limiting"
	MethodGenerateRecoveringFlag      string = "recovering"
	MethodGenerateTracingFlag         string = "tracing"
	MethodGenerateValidatingFlag      string = "validating"
	MethodGenerateValidatorsFlag      string = "validators"
)

const (
	ServiceGenerateGroupGRPC    string = "grpc"
	ServiceGenerateGroupHTTP    string = "http"
	ServiceGenerateGroupMetrics string = "metrics"
)

const (
	MethodGenerateGroupGRPC    string = "grpc"
	MethodGenerateGroupHTTP    string = "http"
	MethodGenerateGroupMetrics string = "metrics"
)
