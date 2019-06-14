package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateProtoBufServiceDefinitionsFile(base string, path string, name string, service *types.Service) *files.ProtoFile {
	file := files.NewProtoFile(base, path, name, true, false)
	file.Pkg = strings.ToLower(strings.ToSnakeCase(service.Name))
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		if len(method.Arguments) > 0 {
			generateProtoBufMethodRequestDefinition(file, method)
		}
		if len(method.Results) > 0 {
			generateProtoBufMethodResponseDefinition(file, method)
		}
	}
	generateProtoBufServiceDefinition(file, service)
	return file
}

func generateProtoBufMethodRequestDefinition(file *files.ProtoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("message %sRequest {", methodName)
	for i, arg := range method.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		file.Pf("\t%s %s = %d", arg.Type.ProtoBufType(), argName, i+1)
	}
	file.Pf("}")
	file.Pf("")
}

func generateProtoBufMethodResponseDefinition(file *files.ProtoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("message %sResponse {", methodName)
	for i, field := range method.Results {
		fieldName := strings.ToUpperFirst(field.Name)
		file.Pf("\t%s %s = %d", field.Type.ProtoBufType(), fieldName, i+1)
	}
	file.Pf("}")
	file.Pf("")
}

func generateProtoBufServiceDefinition(file *files.ProtoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("service %sService {", serviceName)
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		if !method.Options.Generate.GRPCClient && !method.Options.Generate.GRPCServer {
			continue
		}
		if len(method.Arguments) > 0 && len(method.Arguments) > 0 {
			file.Pf("\trpc %s (%sRequest) returns (%sResponse);", methodName, methodName, methodName)
		} else if len(method.Arguments) > 0 {
			file.Pf("\trpc %s (%sRequest) returns (empty.Empty);", methodName, methodName)
		} else if len(method.Results) > 0 {
			file.Pf("\trpc %s (empty.Empty) returns (%sResponse);", methodName, methodName)
		} else {
			file.Pf("\trpc %s (empty.Empty) returns (empty.Empty);", methodName)
		}
	}
	file.Pf("}")
	file.Pf("")
}
