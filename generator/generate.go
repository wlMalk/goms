package generator

import (
	"path/filepath"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/generate_service"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateService(s *types.Service) (files files.Files, err error) {
	serviceNameSnake := strings.ToLower(strings.ToSnakeCase(s.Name))

	files = append(files, generate_service.GenerateRequestsFile(s.Path, filepath.Join("pkg", "service", "requests"), "requests.goms", s.Methods))
	files = append(files, generate_service.GenerateResponseFile(s.Path, filepath.Join("pkg", "service", "responses"), "responses.goms", s.Methods))
	files = append(files, generate_service.GenerateHandlersFile(s.Path, filepath.Join("pkg", "service", "handlers"), "handlers.goms", s))
	files = append(files, generate_service.GenerateConvertersFile(s.Path, filepath.Join("pkg", "service", "handlers", "converters"), "converters.goms", s))
	files = append(files, generate_service.GenerateServiceTransportEndpointsFile(s.Path, filepath.Join("pkg", "transport"), "endpoints.goms", s))

	if helpers.IsMiddlewareEnabled(s) || s.Options.Generate.Middleware {
		files = append(files, generate_service.GenerateServiceMiddlewareFile(s.Path, filepath.Join("pkg", "service", "middleware"), "middleware.goms", s))
	}

	if helpers.IsMiddlewareEnabled(s) && helpers.IsRecoveringEnabled(s) {
		files = append(files, generate_service.GenerateRecoveringMiddlewareFile(s.Path, filepath.Join("pkg", "service", "middleware"), "recovering_middleware.goms", s))
	}

	if helpers.IsMiddlewareEnabled(s) && helpers.IsLoggingEnabled(s) {
		files = append(files, generate_service.GenerateLoggingMiddlewareFile(s.Path, filepath.Join("pkg", "service", "middleware"), "logging_middleware.goms", s))
	}

	if helpers.IsMiddlewareEnabled(s) && helpers.IsCachingEnabled(s) && helpers.IsCachaeble(s) {
		files = append(files, generate_service.GenerateCachingMiddlewareFile(s.Path, filepath.Join("pkg", "service", "middleware"), "caching_middleware.goms", s))
		files = append(files, generate_service.GenerateCachingKeyerFile(s.Path, serviceNameSnake, "caching_keyer", s))
	}

	if helpers.IsMiddlewareEnabled(s) && helpers.IsValidatingEnabled(s) && helpers.IsValidatable(s) {
		files = append(files, generate_service.GenerateValidatingMiddlewareFile(s.Path, filepath.Join("pkg", "service", "middleware"), "validating_middleware.goms", s))
		files = append(files, generate_service.GenerateServiceImplementationValidatorFile(s.Path, serviceNameSnake, "validator", s))
	}

	if helpers.IsMiddlewareEnabled(s) {
		files = append(files, generate_service.GenerateServiceImplementationMiddlewareFile(s.Path, serviceNameSnake, "middleware", s))
	}

	if helpers.IsMethodStubsEnabled(s) {
		files = append(files, generate_service.GenerateServiceImplementationFile(s.Path, serviceNameSnake, serviceNameSnake, s))
	}

	if helpers.IsHTTPServerEnabled(s) {
		files = append(files, generate_service.GenerateHTTPServerFile(s.Path, filepath.Join("pkg", "transport", "http", "server"), "server.goms", s))
	}
	if helpers.IsHTTPClientEnabled(s) {
		files = append(files, generate_service.GenerateHTTPClientFile(s.Path, filepath.Join("pkg", "transport", "http", "client"), "client.goms", s))
		files = append(files, generate_service.GenerateGlobalHTTPClientFile(s.Path, filepath.Join("clients", "http", serviceNameSnake), "client.goms", s))
	}

	if helpers.IsGRPCServerEnabled(s) {
		files = append(files, generate_service.GenerateGRPCServerFile(s.Path, filepath.Join("pkg", "transport", "grpc", "server"), "server.goms", s))
	}
	if helpers.IsGRPCClientEnabled(s) {
		files = append(files, generate_service.GenerateGRPCClientFile(s.Path, filepath.Join("pkg", "transport", "grpc", "client"), "client.goms", s))
		files = append(files, generate_service.GenerateGlobalGRPCClientFile(s.Path, filepath.Join("clients", "grpc", serviceNameSnake), "client.goms", s))
	}

	if helpers.IsHTTPServerEnabled(s) || helpers.IsHTTPClientEnabled(s) {
		files = append(files, generate_service.GenerateHTTPRequestsFile(s.Path, filepath.Join("pkg", "transport", "http", "requests"), "requests.goms", s))
		files = append(files, generate_service.GenerateHTTPResponsesFile(s.Path, filepath.Join("pkg", "transport", "http", "responses"), "responses.goms", s))
		files = append(files, generate_service.GenerateHTTPDecodersFile(s.Path, filepath.Join("pkg", "transport", "http"), "decoders.goms", s))
		files = append(files, generate_service.GenerateHTTPEncodersFile(s.Path, filepath.Join("pkg", "transport", "http"), "encoders.goms", s))
	}

	if helpers.IsGRPCServerEnabled(s) || helpers.IsGRPCClientEnabled(s) {
		files = append(files, generate_service.GenerateGRPCDecodersFile(s.Path, filepath.Join("pkg", "transport", "grpc"), "decoders.goms", s))
		files = append(files, generate_service.GenerateGRPCEncodersFile(s.Path, filepath.Join("pkg", "transport", "grpc"), "encoders.goms", s))
	}

	if s.Options.Generate.ProtoBuf && (helpers.IsGRPCServerEnabled(s) || helpers.IsGRPCClientEnabled(s)) {
		files = append(files, generate_service.GenerateProtoBufServiceDefinitionsFile(s.Path, "proto", "service.goms", s))
		files = append(files, generate_service.GenerateProtoRequestsConvertersFile(s.Path, filepath.Join("pkg", "protobuf", serviceNameSnake, "requests"), "requests.goms", s))
		files = append(files, generate_service.GenerateProtoResponsesConvertersFile(s.Path, filepath.Join("pkg", "protobuf", serviceNameSnake, "responses"), "responses.goms", s))

	}
	if s.Options.Generate.Main {
		files = append(files, generate_service.GenerateServiceMainFile(s.Path, "", "main", s))
		if helpers.IsServerEnabled(s) {
			files = append(files, generate_service.GenerateServiceStartCMDFile(s.Path, filepath.Join("cmd", "start"), "start", s))
		}
	}

	if s.Options.Generate.Dockerfile {
		files = append(files, generate_service.GenerateDockerFile(s.Path, s))
	}

	return
}
