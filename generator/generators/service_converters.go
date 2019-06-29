package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func HandlerConverterTypes(file file.File, service types.Service) error {
	helpers.AddTypesImports(file, service)
	file.Pf("type (")
	file.Pf("requestResponseHandler struct {")
	file.Pf("handler handlers.Handler")
	file.Pf("}")
	file.Pf("endpointHandler struct {")
	file.Pf("handler handlers.RequestResponseHandler")
	file.Pf("}")
	file.Pf(")")
	file.Pf("")
	return nil
}

func HandlerConverterNewFuncs(file file.File, service types.Service) error {
	file.Pf("func HandlerToRequestResponseHandler(h handlers.Handler) handlers.RequestResponseHandler {")
	file.Pf("return &requestResponseHandler{handler: h}")
	file.Pf("}")
	file.Pf("")
	file.Pf("func RequestResponseHandlerToEndpointHandler(h handlers.RequestResponseHandler) handlers.EndpointHandler {")
	file.Pf("return &endpointHandler{handler: h}")
	file.Pf("}")
	file.Pf("")
	return nil
}

func HandlerToRequestResponseHandlerConverter(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		// Import requests
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	argsInCall := append([]string{"ctx"}, helpers.GetMethodArgumentsFromRequestInCall(method.Arguments)...)
	resultVars := append(helpers.GetResultsVarsFromResponse(method.Results), "err")
	file.Pf("func %sHandlerTo%sRequestResponseHandler(next handlers.%sHandler) handlers.%sRequestResponseHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sRequestResponseHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if len(method.Results) > 0 {
		file.Pf("res = &responses.%sResponse{}", methodName)
	}
	file.Pf("%s = next.%s(%s)", strs.Join(resultVars, ", "), methodName, strs.Join(argsInCall, ", "))
	file.Pf("return")
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
	file.Pf("func (h *requestResponseHandler) %s(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if len(method.Results) > 0 {
		file.Pf("res = &responses.%sResponse{}", methodName)
	}
	file.Pf("%s = h.handler.%s(%s)", strs.Join(resultVars, ", "), methodName, strs.Join(argsInCall, ", "))
	file.Pf("return")
	file.Pf("}")
	file.Pf("")
	return nil
}

func RequestResponseHandlerToHandlerConverter(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	argsInCall := helpers.GetMethodArgumentsInCall(method.Arguments)
	returnValues := []string{"err"}
	if len(method.Results) > 0 {
		returnValues = append([]string{"res"}, returnValues...)
	}
	resultVars := append(helpers.GetResultsVarsFromResponse(method.Results), "nil")
	file.Pf("func %sRequestResponseHandlerTo%sHandler(next handlers.%sRequestResponseHandler) handlers.%sHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if len(method.Arguments) > 0 {
		file.Pf("req := requests.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	}
	argsInCall = []string{"ctx"}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		argsInCall = append(argsInCall, "req")
	}
	if len(method.Results) > 0 {
		file.Pf("%s := next.%s(%s)", strs.Join(returnValues, ", "), methodName, strs.Join(argsInCall, ", "))
	} else {
		file.Pf("%s = next.%s(%s)", strs.Join(returnValues, ", "), methodName, strs.Join(argsInCall, ", "))
	}
	file.Pf("if err != nil {")
	file.Pf("return")
	file.Pf("}")
	file.Pf("return %s", strs.Join(resultVars, ", "))
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
	return nil
}

func RequestResponseHandlerToEndpointConverter(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func %sRequestResponseHandlerToEndpoint(next handlers.%sRequestResponseHandler) endpoint.Endpoint {", methodName, methodName)
	file.Pf("return endpoint.Endpoint(func(ctx context.Context, req interface{}) (res interface{}, err error) {")
	retValue := ""
	if len(method.Results) == 0 {
		retValue = "nil, "
	}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		file.Pf("return %snext.%s(ctx, req.(*requests.%sRequest))", retValue, methodName, methodName)
	} else {
		file.Pf("return %snext.%s(ctx)", retValue, methodName)
	}
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
	file.Pf("func (h *endpointHandler) %s(ctx context.Context, req interface{}) (res interface{}, err error) {", methodName)
	retValue = ""
	if len(method.Results) == 0 {
		retValue = "nil, "
	}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		file.Pf("return %sh.handler.%s(ctx, req.(*requests.%sRequest))", retValue, methodName, methodName)
	} else {
		file.Pf("return %sh.handler.%s(ctx)", retValue, methodName)
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func EndpointToRequestResponseHandlerConverter(file file.File, service types.Service, method types.Method) error {
	file.AddImport("", "context")
	file.AddImport("", service.ImportPath, "/pkg/service/handlers")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", service.ImportPath, "/pkg/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func EndpointTo%sRequestResponseHandler(next endpoint.Endpoint) handlers.%sRequestResponseHandler {", methodName, methodName)
	file.Pf("return handlers.%sRequestResponseHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	req := "nil"
	if len(method.Arguments) > 0 {
		req = "req"
	}
	if len(method.Results) > 0 {
		file.Pf("resp, err := next(ctx, %s)", req)
	} else {
		file.Pf("_, err = next(ctx, %s)", req)
	}
	if len(method.Results) > 0 {
		file.Pf("if err != nil {")
		file.Pf("return nil, err")
		file.Pf("}")
		file.Pf("if resp != nil {")
		file.Pf("return resp.(*responses.%sResponse), nil", methodName)
		file.Pf("}")
	}
	file.Pf("return")
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
	return nil
}
