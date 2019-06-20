package generators

import (
	"github.com/wlMalk/goms/constants"
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ProtoBufPackageDefinition(file file.File, service types.Service) error {
	file.Pf("option go_package = \"%s/pkg/protobuf/%s\";", service.ImportPath, strings.ToLower(strings.ToSnakeCase(service.Name)))
	file.P("")
	return nil
}

func ProtoBufMethodRequestDefinition(file file.File, service types.Service, method types.Method) error {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("message %sRequest {", methodName)
	for i, arg := range method.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		file.Pf("\t%s %s = %d;", arg.Type.ProtoBufType(), argName, i+1)
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func ProtoBufMethodResponseDefinition(file file.File, service types.Service, method types.Method) error {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("message %sResponse {", methodName)
	for i, field := range method.Results {
		fieldName := strings.ToUpperFirst(field.Name)
		file.Pf("\t%s %s = %d;", field.Type.ProtoBufType(), fieldName, i+1)
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func ProtoBufServiceDefinition(file file.File, service types.Service) error {
	serviceName := strings.ToUpperFirst(service.Name)
	file.Pf("service %sService {", serviceName)
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		if method.Generate.HasNone(constants.MethodGenerateGRPCServerFlag, constants.MethodGenerateGRPCClientFlag) {
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
	return nil
}
