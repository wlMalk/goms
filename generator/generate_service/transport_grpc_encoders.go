package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateGRPCEncodersFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		generateGRPCRequestEncoder(file, method)
		generateGRPCResponseEncoder(file, method)
	}
	return file
}

func generateGRPCRequestEncoder(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(method.Service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sRequest(ctx context.Context, request interface{}) (interface{}, error) {", methodName)
	if len(method.Arguments) > 0 {
		file.AddImport("proto_requests", method.Service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(method.Service.Name)), "requests")
		file.AddImport("", method.Service.ImportPath, "/pkg/service/requests")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if request == nil {")
		file.Pf("return nil, errors.InvalidRequest(\"%s\", \"%s\")", helpers.GetName(serviceName, method.Service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return proto_requests.%s(request.(*requests.%sRequest))", methodName, methodName)
	} else {
		file.AddImport("", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("return &empty.Empty{}, nil")
	}
	file.Pf("}")
	file.Pf("")
}

func generateGRPCResponseEncoder(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(method.Service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sResponse(ctx context.Context, response interface{}) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("proto_responses", method.Service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(method.Service.Name)), "responses")
		file.AddImport("", method.Service.ImportPath, "/pkg/service/responses")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if response == nil {")
		file.Pf("return nil, errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, method.Service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return proto_responses.%s(response.(*responses.%sResponse))", methodName, methodName)
	} else {
		file.AddImport("", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("return &empty.Empty{}, nil")
	}
	file.Pf("}")
	file.Pf("")
}
