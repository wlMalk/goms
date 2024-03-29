package generators

import (
	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GRPCTransportServerRegisterFunc(file file.File, service types.Service) error {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("kit_grpc", "github.com/go-kit/kit/transport/grpc")
	file.AddImport("", service.ImportPath, "/pkg/transport")
	file.AddImport(serviceNameSnake+"_grpc", service.ImportPath, "/pkg/transport/grpc")
	file.AddImport("goms_grpc", "github.com/wlMalk/goms/goms/transport/grpc")
	file.Pf("func Register(server *goms_grpc.Server, endpoints *transport.%s, opts ...kit_grpc.ServerOption) {", serviceName)
	file.Pf("RegisterSpecial(server, endpoints, func(_ string) []kit_grpc.ServerOption {")
	file.Pf("return opts")
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
	return nil
}

func GRPCTransportServerRegisterSpecialFunc(file file.File, service types.Service) error {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("kit_grpc", "github.com/go-kit/kit/transport/grpc")
	file.AddImport("", service.ImportPath, "/pkg/transport")
	file.AddImport(serviceNameSnake+"_grpc", service.ImportPath, "/pkg/transport/grpc")
	file.AddImport("goms_grpc", "github.com/wlMalk/goms/goms/transport/grpc")
	file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
	file.Pf("func RegisterSpecial(server *goms_grpc.Server, endpoints *transport.%s, optionsFunc func(method string) (opts []kit_grpc.ServerOption)) {", serviceName)
	file.Pf("handler := &serverHandler{")
	for _, method := range helpers.GetMethodsWithGRPCServerEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: kit_grpc.NewServer(", lowerMethodName)
		file.Pf("endpoints.%s,", methodName)
		file.Pf("%s_grpc.Decode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_grpc.Encode%sResponse,", serviceNameSnake, methodName)
		file.Pf("optionsFunc(\"%s\")...),", helpers.GetName(methodName, method.Alias))
	}
	file.Pf("}")
	file.Pf("pb.Register%sServiceServer(server.Server, handler)", serviceName)
	file.Pf("}")
	file.Pf("")
	return nil
}

func GRPCTransportServerHandlerStruct(file file.File, service types.Service) error {
	file.Pf("type serverHandler struct {")
	for _, method := range helpers.GetMethodsWithGRPCServerEnabled(service) {
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s kit_grpc.Handler", lowerMethodName)
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func GRPCTransportServerHandlerMethodFunc(file file.File, service types.Service, method types.Method) error {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	file.AddImport("", "context")
	if len(method.Arguments) > 0 && len(method.Results) > 0 {
		file.AddImport("pb", service.ImportPath, "pkg/protobuf", strings.ToLower(strings.ToSnakeCase(service.Name)))
		file.Pf("func (h *serverHandler) %s(ctx context.Context, req *pb.%sRequest) (*pb.%sResponse, error) {", methodName, methodName, methodName)
	} else if len(method.Arguments) > 0 {
		file.AddImport("empty", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("func (h *serverHandler) %s(ctx context.Context, req *pb.%sRequest) (*empty.Empty, error) {", methodName, methodName)
	} else if len(method.Results) > 0 {
		file.AddImport("empty", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("func (h *serverHandler) %s(ctx context.Context, req *empty.Empty) (*pb.%sResponse, error) {", methodName, methodName)
	} else {
		file.AddImport("empty", "github.com/golang/protobuf/ptypes/empty")
		file.Pf("func (h *serverHandler) %s(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {", methodName)
	}
	file.Pf("_, resp, err := h.%s.ServeGRPC(ctx, req)", lowerMethodName)
	file.Pf("if err != nil {")
	file.Pf("return nil, err")
	file.Pf("}")
	if len(method.Results) > 0 {
		file.Pf("return resp.(*pb.%sResponse), nil", methodName)
	} else {
		file.Pf("return resp.(*empty.Empty), nil")
	}
	file.Pf("}")
	file.Pf("")
	return nil
}
