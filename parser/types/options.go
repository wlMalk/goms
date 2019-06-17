package types

type TagOptions map[string]interface{}

type TagsOptions map[string]TagOptions

type ServiceOptions struct {
	HTTP     HTTPServiceOptions
	GRPC     GRPCServiceOptions
	Generate GenerateServiceOptions
}

type HTTPServiceOptions struct {
	URIPrefix string
}

type GRPCServiceOptions struct {
}

type GenerateServiceOptions struct {
	Logger           bool
	CircuitBreaking  bool
	RateLimiting     bool
	Recovering       bool
	Caching          bool
	Logging          bool
	Tracing          bool
	ServiceDiscovery bool
	ProtoBuf         bool
	Main             bool
	Validators       bool
	Validating       bool
	Middleware       bool
	MethodStubs      bool
	FrequencyMetric  bool
	LatencyMetric    bool
	CounterMetric    bool
	HTTPServer       bool
	HTTPClient       bool
	GRPCServer       bool
	GRPCClient       bool
	Dockerfile       bool
}

type MethodOptions struct {
	HTTP     HTTPMethodOptions
	GRPC     GRPCMethodOptions
	Logging  LoggingMethodOptions
	Generate GenerateMethodOptions
}

type GenerateMethodOptions struct {
	CircuitBreaking bool
	RateLimiting    bool
	Recovering      bool
	Caching         bool
	Logging         bool
	Validators      bool
	Validating      bool
	Middleware      bool
	MethodStubs     bool
	Tracing         bool
	FrequencyMetric bool
	LatencyMetric   bool
	CounterMetric   bool
	HTTPServer      bool
	HTTPClient      bool
	GRPCServer      bool
	GRPCClient      bool
}

type HTTPMethodOptions struct {
	Method string
	URI    string
	AbsURI string
}

type GRPCMethodOptions struct {
}

type LoggingMethodOptions struct {
	IgnoredArguments []string
	IgnoredResults   []string
	LenArguments     []string
	LenResults       []string
	IgnoreError      bool
}

type ArgumentOptions struct {
	HTTP HTTPArgumentOptions
}

type HTTPArgumentOptions struct {
	Origin string
}
