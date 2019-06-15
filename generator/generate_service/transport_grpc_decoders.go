package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateGRPCDecodersFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		generateGRPCRequestDecoder(file, method)
		generateGRPCResponseDecoder(file, method)
	}
	return file
}

func generateGRPCRequestDecoder(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(method.Service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sRequest(ctx context.Context, req interface{}) (interface{}, error) {", methodName)
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(method.Service.Name)), "requests")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.AddImport("pb", method.Service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(method.Service.Name)))
		file.Pf("if req == nil {")
		file.Pf("return nil, errors.InvalidRequest(\"%s\", \"%s\")", helpers.GetName(serviceName, method.Service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return requests.%sFromProto(req.(*pb.%sRequest))", methodName, methodName)
	} else {
		file.Pf("return nil, nil")
	}
	file.Pf("}")
	file.Pf("")
}

func generateGRPCResponseDecoder(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	serviceName := strings.ToUpperFirst(method.Service.Name)
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func Decode%sResponse(ctx context.Context, res interface{}) (interface{}, error) {", methodName)
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(method.Service.Name)), "responses")
		file.AddImport("", "github.com/wlMalk/goms/goms/errors")
		file.AddImport("pb", method.Service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(method.Service.Name)))
		file.Pf("if res == nil {")
		file.Pf("return nil, errors.InvalidResponse(\"%s\", \"%s\")", helpers.GetName(serviceName, method.Service.Alias), helpers.GetName(methodName, method.Alias))
		file.Pf("}")
		file.Pf("return responses.%sFromProto(res.(*pb.%sResponse))", methodName, methodName)
	} else {
		file.Pf("return nil, nil")
	}
	file.Pf("}")
	file.Pf("")
}
