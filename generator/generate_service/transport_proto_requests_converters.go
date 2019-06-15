package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateProtoRequestsConvertersFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		generateProtoRequestNewFunc(file, method)
		generateProtoRequestNewProtoFunc(file, method)
	}
	return file
}

func generateProtoRequestNewFunc(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "/pkg/service/requests")
	file.AddImport("pb", method.Service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(method.Service.Name)))
	file.Pf("func %s(r *requests.%sRequest) (req *pb.%sRequest, err error) {", methodName, methodName, methodName)
	file.Pf("req = &pb.%sRequest{}", methodName)
	for _, arg := range method.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		file.Pf("req.%s = r.%s", argName, argName)
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}

func generateProtoRequestNewProtoFunc(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "/pkg/service/requests")
	file.AddImport("pb", method.Service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(method.Service.Name)))
	file.Pf("func %sFromProto(r *pb.%sRequest) (req *requests.%sRequest, err error) {", methodName, methodName, methodName)
	file.Pf("req = &requests.%sRequest{}", methodName)
	for _, arg := range method.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		file.Pf("req.%s = r.%s", argName, argName)
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}
