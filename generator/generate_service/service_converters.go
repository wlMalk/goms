package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateConvertersFile(base string, path string, name string, methods []*types.Method) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range methods {
		generateRequestHandlerToHandlerConverter(file, method)
		generateHandlerToRequestHandlerConverter(file, method)
		generateRequestResponseHandlerToRequestHandlerConverter(file, method)
		generateRequestHandlerToRequestResponseHandlerConverter(file, method)
		generateEndpointToRequestResponseConverter(file, method)
		generateRequestResponseHandlerToEndpointConverter(file, method)
	}
	return file
}

func generateRequestHandlerToHandlerConverter(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	methodName := strings.ToUpperFirst(method.Name)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		// Import requests
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	argsInCall := append([]string{"ctx"}, helpers.GetMethodArgumentsFromRequestInCall(method.Arguments)...)
	file.Pf("func %sRequestHandlerTo%sHandler(next handlers.%sHandler) handlers.%sRequestHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sRequestHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return next.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateHandlerToRequestHandlerConverter(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, helpers.GetMethodArguments(method.Arguments)...)
	results := append(helpers.GetMethodResults(method.Results), "err error")
	argsInCall := helpers.GetMethodArgumentsInCall(method.Arguments)
	file.Pf("func %sHandlerTo%sRequestHandler(next handlers.%sRequestHandler) handlers.%sHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if len(method.Arguments) > 0 {
		file.Pf("req := requests.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	}
	argsInCall = []string{"ctx"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		argsInCall = append(argsInCall, "req")
	}
	file.Pf("return next.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateRequestResponseHandlerToRequestHandlerConverter(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	resultVars := append(helpers.GetResultsVarsFromResponse(method.Results), "err")
	file.Pf("func %sRequestResponseHandlerTo%sRequestHandler(next handlers.%sRequestHandler) handlers.%sRequestResponseHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sRequestResponseHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	if len(method.Results) > 0 {
		file.Pf("res = &responses.%sResponse{}", methodName)
	}
	file.Pf("%s = next.%s(%s)", strs.Join(resultVars, ", "), methodName, strs.Join(argsInCall, ", "))
	file.Pf("return")
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateRequestHandlerToRequestResponseHandlerConverter(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := append(helpers.GetMethodResults(method.Results), "err error")
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	returnValues := []string{"err"}
	if len(method.Results) > 0 {
		returnValues = append([]string{"res"}, returnValues...)
	}
	resultVars := append(helpers.GetResultsVarsFromResponse(method.Results), "nil")
	file.Pf("func %sRequestHandlerTo%sRequestResponseHandler(next handlers.%sRequestResponseHandler) handlers.%sRequestHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sRequestHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
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
}

func generateEndpointToRequestResponseConverter(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func EndpointTo%sRequestResponseHandler(next handlers.%sRequestResponseHandler) endpoint.Endpoint {", methodName, methodName)
	file.Pf("return endpoint.Endpoint(func(ctx context.Context, req interface{}) (res interface{}, err error) {")
	retValue := ""
	if len(method.Results) == 0 {
		retValue = "nil, "
	}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		file.Pf("return %snext.%s(ctx, req.(*requests.%sRequest))", retValue, methodName, methodName)
	} else {
		file.Pf("return %snext.%s(ctx)", retValue, methodName)
	}
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateRequestResponseHandlerToEndpointConverter(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
	file.AddImport("", method.Service.ImportPath, "/service/handlers")
	file.AddImport("", "github.com/go-kit/kit/endpoint")
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/requests")
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		file.AddImport("", method.Service.ImportPath, "/service/responses")
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	file.Pf("func %sRequestResponseHandlerToEndpoint(next endpoint.Endpoint) handlers.%sRequestResponseHandler {", methodName, methodName)
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
}