package generator

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPClientFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateHTTPTransportClientStruct(file, service)
	generateHTTPTransportClientNewFunc(file, service)
	for _, method := range service.Methods {
		generateHTTPTransportClientMethodFunc(file, method)
	}
	generateHTTPTransportClientGlobalVar(file, service)
	for _, method := range service.Methods {
		generateHTTPTransportClientGlobalFunc(file, method)
	}
	return file
}

func generateHTTPTransportClientStruct(file *GoFile, service *types.Service) {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/service/handlers")
	file.Pf("type Client struct {")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s handlers.%sHandler", lowerMethodName, methodName)
	}
	file.Pf("}")
	file.Pf("")
}

func generateHTTPTransportClientNewFunc(file *GoFile, service *types.Service) {
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("", "net/url")
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
	file.AddImport("", service.ImportPath, "/service/handlers/converters")
	file.AddImport(serviceNameSnake+"_http", service.ImportPath, "/service/transport/http")
	file.Pf("func New(u *url.URL, opts ...kit_http.ClientOption) *Client {")
	file.Pf("return &Client{")
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		lowerMethodName := strings.ToLowerFirst(method.Name)
		file.Pf("%s: converters.%sHandlerTo%sRequestHandler(", lowerMethodName, methodName, methodName)
		file.Pf("converters.%sRequestHandlerTo%sRequestResponseHandler(", methodName, methodName)
		file.Pf("converters.%sRequestResponseHandlerToEndpoint(", methodName)
		file.Pf("kit_http.NewClient(")
		file.Pf("\"POST\", u,")
		file.Pf("%s_http.Encode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_http.Decode%sResponse,", serviceNameSnake, methodName)
		file.Pf("opts...,")
		file.Pf(").Endpoint()))),")
	}
	file.Pf("}")
	file.Pf("}")
	file.Pf("")
}

func generateHTTPTransportClientMethodFunc(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	argsInCall := append([]string{"ctx"}, getMethodArgumentsInCall(method.Arguments)...)
	file.Pf("func (c *Client) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return c.%s.%s(%s)", lowerMethodName, methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}

func generateHTTPTransportClientGlobalVar(file *GoFile, service *types.Service) {
	file.Pf("var client *Client = New(nil)")
	file.Pf("")
}

func generateHTTPTransportClientGlobalFunc(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	lowerMethodName := strings.ToLowerFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	argsInCall := append([]string{"ctx"}, getMethodArgumentsInCall(method.Arguments)...)
	file.Pf("func %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return client.%s.%s(%s)", lowerMethodName, methodName, strs.Join(argsInCall, ", "))
	file.Pf("}")
	file.Pf("")
}