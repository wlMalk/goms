package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateProtoResponsesConvertersFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range helpers.GetMethodsWithGRPCEnabled(service) {
		generateProtoResponseNewFunc(file, method)
		generateProtoResponseNewProtoFunc(file, method)
	}
	return file
}

func generateProtoResponseNewFunc(file *files.GoFile, method *types.Method) {
	if len(method.Results) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "/pkg/service/responses")
	file.AddImport("pb", method.Service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(method.Service.Name)))
	file.Pf("func %s(r *responses.%sResponse) (res *pb.%sResponse, err error) {", methodName, methodName, methodName)
	file.Pf("res = &pb.%sResponse{}", methodName)
	for _, res := range method.Results {
		resName := strings.ToUpperFirst(res.Name)
		file.Pf("res.%s = r.%s", resName, resName)
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}

func generateProtoResponseNewProtoFunc(file *files.GoFile, method *types.Method) {
	if len(method.Results) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "/pkg/service/responses")
	file.AddImport("pb", method.Service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(method.Service.Name)))
	file.Pf("func %sFromProto(r *pb.%sResponse) (res *responses.%sResponse, err error) {", methodName, methodName, methodName)
	file.Pf("res = &responses.%sResponse{}", methodName)
	for _, res := range method.Results {
		resName := strings.ToUpperFirst(res.Name)
		file.Pf("res.%s = r.%s", resName, resName)
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
}
