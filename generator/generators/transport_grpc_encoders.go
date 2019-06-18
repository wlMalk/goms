package generators

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GRPCRequestEncoder(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sRequest(ctx context.Context, request interface{}) (interface{}, error) {", methodName)
	if len(method.Arguments) > 0 {
		file.AddImport("proto_requests", service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)), "requests")
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if request == nil {")
		file.Pf("return nil, errors.InvalidRequest(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return proto_requests.%s(request.(*requests.%sRequest))", methodName, methodName)
	} else {
		file.AddImport("", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("return &empty.Empty{}, nil")
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func GRPCResponseEncoder(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Encode%sResponse(ctx context.Context, response interface{}) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("proto_responses", service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)), "responses")
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.Pf("if response == nil {")
		file.Pf("return nil, errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return proto_responses.%s(response.(*responses.%sResponse))", methodName, methodName)
	} else {
		file.AddImport("", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("return &empty.Empty{}, nil")
	}
	file.Pf("}")
	file.Pf("")
	return nil
}
