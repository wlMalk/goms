package generator

import "github.com/wlMalk/goms/generator/file"

var builtInGenerators []GeneratorOption = []GeneratorOption{
	DockerfileFileSpec,
	ProtoBufServiceDefinitionsFileSpec,
	ServiceMainFileSpec,
	ServiceStartCMDFileSpec,
	CachingMiddlewareFileSpec,
	ConvertersFileSpec,
	HandlersFileSpec,
	ServiceImplementationFileSpec,
	ServiceImplementationValidatorFileSpec,
	ServiceImplementationMiddlewareFileSpec,
	CachingKeyerFileSpec,
	LoggingMiddlewareFileSpec,
	ServiceMiddlewareFileSpec,
	RecoveringMiddlewareFileSpec,
	ServiceTransportEndpointsFileSpec,
	ValidatingMiddlewareFileSpec,
	GRPCClientFileSpec,
	GlobalGRPCClientFileSpec,
	GRPCDecodersFileSpec,
	GRPCEncodersFileSpec,
	GRPCServerFileSpec,
	HTTPClientFileSpec,
	GlobalHTTPClientFileSpec,
	HTTPDecodersFileSpec,
	HTTPEncodersFileSpec,
	HTTPRequestsFileSpec,
	HTTPResponsesFileSpec,
	HTTPServerFileSpec,
	ProtoRequestsConvertersFileSpec,
	ProtoResponsesConvertersFileSpec,
	RequestsFileSpec,
	ResponseFileSpec,
	ServiceTypesDefinitionsFileSpec,
}

var builtInFileCreators []GeneratorOption = []GeneratorOption{
	GeneratorOption(func(generator *Generator) {
		generator.AddCreator("go", GoFileCreator)
	}),
	GeneratorOption(func(generator *Generator) {
		generator.AddCreator("proto", ProtoFileCreator)
	}),
	GeneratorOption(func(generator *Generator) {
		generator.AddCreator("Dockerfile", TextFileCreator(""))
	}),
}

func Default(opts ...GeneratorOption) *Generator {
	return New(append(append(builtInGenerators, builtInFileCreators...), opts...)...)
}

func New(opts ...GeneratorOption) *Generator {
	g := &Generator{
		creators: map[string]file.Creator{},
		specs:    map[string]file.Spec{},
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}
