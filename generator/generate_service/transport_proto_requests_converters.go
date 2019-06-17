package generate_service

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ProtoRequestNewFunc(file file.File, service types.Service, method types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
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

func ProtoRequestNewProtoFunc(file file.File, service types.Service, method types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
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
