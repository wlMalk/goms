package generators

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ProtoResponseNewFunc(file file.File, service types.Service, method types.Method) error {
	if len(method.Results) == 0 {
		return nil
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", service.ImportPath, "/pkg/service/responses")
	file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
	file.Pf("func %s(r *responses.%sResponse) (res *pb.%sResponse, err error) {", methodName, methodName, methodName)
	file.Pf("res = &pb.%sResponse{}", methodName)
	for _, res := range method.Results {
		resName := strings.ToUpperFirst(res.Name)
		file.Pf("res.%s = r.%s", resName, resName)
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
	return nil
}

func ProtoResponseNewProtoFunc(file file.File, service types.Service, method types.Method) error {
	if len(method.Results) == 0 {
		return nil
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", service.ImportPath, "/pkg/service/responses")
	file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
	file.Pf("func %sFromProto(r *pb.%sResponse) (res *responses.%sResponse, err error) {", methodName, methodName, methodName)
	file.Pf("res = &responses.%sResponse{}", methodName)
	for _, res := range method.Results {
		resName := strings.ToUpperFirst(res.Name)
		file.Pf("res.%s = r.%s", resName, resName)
	}
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
	return nil
}
