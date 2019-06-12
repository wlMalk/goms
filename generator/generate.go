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

	files = append(files, generate_service.GenerateRequestsFile(s.Path, filepath.Join("service", "requests"), "requests.goms", s.Methods))
	files = append(files, generate_service.GenerateResponseFile(s.Path, filepath.Join("service", "responses"), "responses.goms", s.Methods))
	files = append(files, generate_service.GenerateHandlersFile(s.Path, filepath.Join("service", "handlers"), "handlers.goms", s))
	files = append(files, generate_service.GenerateConvertersFile(s.Path, filepath.Join("service", "handlers", "converters"), "converters.goms", s))
	files = append(files, generate_service.GenerateServiceTransportEndpointsFile(s.Path, filepath.Join("service", "transport"), "endpoints.goms", s))

	if helpers.IsMiddlewareEnabled(s) || helpers.IsLoggingEnabled(s) || helpers.IsCachingEnabled(s) {
		files = append(files, generate_service.GenerateServiceMiddlewareFile(s.Path, filepath.Join("service", "middleware"), "middleware.goms", s))
	}

	if helpers.IsLoggingEnabled(s) {
		files = append(files, generate_service.GenerateLoggingMiddlewareFile(s.Path, filepath.Join("service", "middleware"), "logging_middleware.goms", s))
	}

	if helpers.IsCachingEnabled(s) {
		files = append(files, generate_service.GenerateCachingMiddlewareFile(s.Path, filepath.Join("service", "middleware"), "caching_middleware.goms", s))
		files = append(files, generate_service.GenerateCachingKeyerFile(s.Path, strings.ToLowerFirst(s.Name), "caching_keyer", s))
	}

	if helpers.IsMethodStubsEnabled(s) {
		files = append(files, generate_service.GenerateServiceImplementationFile(s.Path, strings.ToLowerFirst(s.Name), strings.ToLowerFirst(s.Name), s))
	}

	if helpers.IsValidatorsEnabled(s) {
		files = append(files, generate_service.GenerateServiceImplementationValidatorsFile(s.Path, strings.ToLowerFirst(s.Name), "validators", s))
	}

	if helpers.IsMiddlewareEnabled(s) {
		files = append(files, generate_service.GenerateServiceImplementationMiddlewareFile(s.Path, strings.ToLowerFirst(s.Name), "middleware", s))
	}

	if helpers.IsHTTPServerEnabled(s) {
		files = append(files, generate_service.GenerateHTTPServerFile(s.Path, filepath.Join("service", "transport", "http", "server"), "server.goms", s))
	}
	if helpers.IsHTTPClientEnabled(s) {
		files = append(files, generate_service.GenerateHTTPClientFile(s.Path, filepath.Join("service", "transport", "http", "client"), "client.goms", s))
	}

	if helpers.IsHTTPServerEnabled(s) || helpers.IsHTTPClientEnabled(s) {
		files = append(files, generate_service.GenerateHTTPRequestsFile(s.Path, filepath.Join("service", "transport", "http", "requests"), "requests.goms", s))
		files = append(files, generate_service.GenerateHTTPResponsesFile(s.Path, filepath.Join("service", "transport", "http", "responses"), "responses.goms", s))
		files = append(files, generate_service.GenerateHTTPDecodersFile(s.Path, filepath.Join("service", "transport", "http"), "decoders.goms", s))
		files = append(files, generate_service.GenerateHTTPEncodersFile(s.Path, filepath.Join("service", "transport", "http"), "encoders.goms", s))
	}

	if s.Options.Generate.Main {
		files = append(files, generate_service.GenerateServiceMainFile(s.Path, filepath.Join("cmd", strings.ToLowerFirst(s.Name)), "main", s))
	}
	return
}
