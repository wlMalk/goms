package generator

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateRequestHandlerToHandlerConverter(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	results := append(getMethodResults(method.Results), "err error")
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		// Import requests
		args = append(args, "req *requests."+methodName+"Request")
	}
	argsInCall := append([]string{"ctx"}, getMethodArgumentsFromRequestInCall(method.Arguments)...)
	file.Pf("func %sRequestHandlerTo%sHandler(next handlers.%sHandler) handlers.%sRequestHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sRequestHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("return next.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateHandlerToRequestHandlerConverter(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	args := append([]string{"ctx context.Context"}, getMethodArguments(method.Arguments)...)
	results := append(getMethodResults(method.Results), "err error")
	argsInCall := getMethodArgumentsInCall(method.Arguments)
	file.Pf("func %sHandlerTo%sRequestHandler(next handlers.%sRequestHandler) handlers.%sHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("req := requests.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	argsInCall = []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	file.Pf("return next.%s(%s)", methodName, strs.Join(argsInCall, ", "))
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateRequestResponseHandlerToRequestHandlerConverter(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		// Import requests
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := []string{"err error"}
	if len(method.Results) > 0 {
		// Import responses
		results = append([]string{"res *responses." + methodName + "Response"}, results...)
	}
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	resultVars := append(getResultsVarsFromResponse(method.Results), "err")
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

func generateRequestHandlerToRequestResponseHandlerConverter(file *GoFile, method *types.Method) {
	methodName := strings.ToUpperFirst(method.Name)
	args := []string{"ctx context.Context"}
	if len(method.Arguments) > 0 {
		// Import requests
		args = append(args, "req *requests."+methodName+"Request")
	}
	results := append(getMethodResults(method.Results), "err error")
	argsInCall := []string{"ctx"}
	if len(method.Arguments) > 0 {
		argsInCall = append(argsInCall, "req")
	}
	returnValues := []string{"err"}
	if len(method.Results) > 0 {
		returnValues = append([]string{"res"}, returnValues...)
	}
	resultVars := append(getResultsVarsFromResponse(method.Results), "nil")
	file.Pf("func %sRequestHandlerTo%sRequestResponseHandler(next handlers.%sRequestResponseHandler) handlers.%sRequestHandler {", methodName, methodName, methodName, methodName)
	file.Pf("return handlers.%sRequestHandlerFunc(func(%s) (%s) {", methodName, strs.Join(args, ", "), strs.Join(results, ", "))
	file.Pf("%s := next.%s(%s)", strs.Join(returnValues, ", "), methodName, strs.Join(argsInCall, ", "))
	file.Pf("if err != nil {")
	file.Pf("return")
	file.Pf("}")
	file.Pf("return %s", strs.Join(resultVars, ", "))
	file.Pf("})")
	file.Pf("}")
	file.Pf("")
}

func generateConvertersFile(base string, path string, name string, methods []*types.Method) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	for _, method := range methods {
		generateRequestHandlerToHandlerConverter(file, method)
		generateHandlerToRequestHandlerConverter(file, method)
		generateRequestResponseHandlerToRequestHandlerConverter(file, method)
		generateRequestHandlerToRequestResponseHandlerConverter(file, method)
	}
	return file
}
