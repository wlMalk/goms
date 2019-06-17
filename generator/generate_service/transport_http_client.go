package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func HTTPTransportClientStruct(file file.File, service types.Service) {
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

func HTTPTransportClientNewFunc(file file.File, service types.Service) {
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("", "net/url")
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers/converters")
	file.AddImport(serviceNameSnake+"_http", service.ImportPath, "/pkg/transport/http")
	file.Pf("func New(u *url.URL, opts ...kit_http.ClientOption) *Client {")
	file.Pf("return &Client{")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: converters.%sRequestResponseHandlerTo%sHandler(", lowerMethodName, methodName, methodName)
		file.Pf("converters.EndpointTo%sRequestResponseHandler(", methodName)
		file.Pf("kit_http.NewClient(")
		file.Pf("\"POST\", u,")
		file.Pf("%s_http.Encode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_http.Decode%sResponse,", serviceNameSnake, methodName)
		file.Pf("opts...,")
		file.Pf(").Endpoint())),")
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}

func HTTPTransportClientMethodFunc(file file.File, service types.Service, method types.Method) {
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

func HTTPTransportClientGlobalVar(file file.File, service types.Service) {
	file.AddImport("", service.ImportPath, "/pkg/transport/http/client")
	file.Pf("var c *client.Client = client.New(nil)")
	file.Pf("")
}

func HTTPTransportClientGlobalFunc(file file.File, service types.Service, method types.Method) {
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
