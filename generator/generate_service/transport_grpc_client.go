package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GRPCTransportClientStruct(file file.File, service types.Service) {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.Pf("type Client struct {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s handlers.%sHandler", lowerMethodName, methodName)
	}
	file.Pf("}")
	file.Pf("")
}

func GRPCTransportClientNewFunc(file file.File, service types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("kit_grpc", "github.com/go-kit/kit/transport/grpc")
	file.AddImport("", "google.golang.org/grpc")
	file.Pf("func New(conn *grpc.ClientConn, opts ...kit_grpc.ClientOption) *Client {")
	file.Pf("return &Client{")
	for _, method := range helpers.GetMethodsWithGRPCClientEnabled(service) {
		file.AddImport("", service.ImportPath, "/pkg/service/handlers/converters")
		file.AddImport(serviceNameSnake+"_grpc", service.ImportPath, "/pkg/transport/grpc")
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: converters.%sRequestResponseHandlerTo%sHandler(", lowerMethodName, methodName, methodName)
		file.Pf("converters.EndpointTo%sRequestResponseHandler(", methodName)
		file.Pf("kit_grpc.NewClient(")
		file.Pf("conn, \"%s.%sService\", \"%s\",", serviceNameSnake, serviceName, methodName)
		file.Pf("%s_grpc.Encode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_grpc.Decode%sResponse,", serviceNameSnake, methodName)
		if len(method.Results) > 0 {
			file.AddImport("pb", service.ImportPath, "/pkg/protobuf/", strings.ToLower(strings.ToSnakeCase(service.Name)))
			file.Pf("pb.%sResponse{},", methodName)
		} else {
			file.AddImport("", "github.com/golang/protobuf/ptypes/empty")
			file.Pf("empty.Empty{},")
		}
		file.Pf("opts...,")
		file.Pf(").Endpoint())),")
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}

func GRPCTransportClientMethodFunc(file file.File, service types.Service, method types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	argsInCall := append([]string{"ctx"}, helpers.GetMethodArgumentsInCall(method.Arguments)...)
	file.Pf("func (c *Client) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return c.%s.%s(%s)", lowerMethodName, methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}

func GRPCTransportClientGlobalVar(file file.File, service types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/transport/grpc/client")
	file.Pf("var c *client.Client = client.New(nil)")
	file.Pf("")
}

func GRPCTransportClientGlobalFunc(file file.File, service types.Service, method types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", "context")
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	argsInCall := append([]string{"ctx"}, helpers.GetMethodArgumentsInCall(method.Arguments)...)
	file.Pf("func %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return c.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}
