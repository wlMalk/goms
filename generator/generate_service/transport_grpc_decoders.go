package generate_service

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GRPCRequestDecoder(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sRequest(ctx context.Context, req interface{}) (interface{}, error) {", methodName)
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)), "requests")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
		file.Pf("if req == nil {")
		file.Pf("return nil, errors.InvalidRequest(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return requests.%sFromProto(req.(*pb.%sRequest))", methodName, methodName)
	} else {
		file.Pf("return nil, nil")
	}
	file.Pf("}")
	file.Pf("")
}

func GRPCResponseDecoder(file file.File, service types.Service, method types.Method) {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sResponse(ctx context.Context, res interface{}) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)), "responses")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
		file.Pf("if res == nil {")
		file.Pf("return nil, errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return responses.%sFromProto(res.(*pb.%sResponse))", methodName, methodName)
	} else {
		file.Pf("return nil, nil")
	}
	file.Pf("}")
	file.Pf("")
}
