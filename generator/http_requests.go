package generator

import (
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPRequestsFile(base string, path string, name string, methods []*types.Method) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	for _, method := range methods {
		generateHTTPRequest(file, method)
	}
	return file
}

func generateHTTPRequest(file *GoFile, method *types.Method) {
	if (method.Options.HTTP.Method == "POST" || method.Options.HTTP.Method == "PUT") && hasArgumentsOfOrigin(method.Arguments, "BODY") {
		file.Pf("type " + method.Name + "RequestBody struct {")
		generateHTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "BODY"))
		file.Pf("}")
		file.Pf("")
	}
	file.Pf("type " + method.Name + "Request struct {")
	if (method.Options.HTTP.Method == "POST" || method.Options.HTTP.Method == "PUT") && hasArgumentsOfOrigin(method.Arguments, "BODY") {
		file.Pf("Body *" + method.Name + "RequestBody")
		file.Pf("")
	}
	if hasArgumentsOfOrigin(method.Arguments, "HEADER") {
		generateHTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "HEADER"))
	}
	if hasArgumentsOfOrigin(method.Arguments, "QUERY") {
		generateHTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "QUERY"))
	}
	if hasArgumentsOfOrigin(method.Arguments, "PATH") {
		generateHTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "PATH"))
	}
	file.Pf("}")
	file.Pf("")
}

func hasArgumentsOfOrigin(args []*types.Argument, origin string) bool {
	for _, arg := range args {
		if arg.Options.HTTP.Origin == origin {
			return true
		}
	}
	return false
}

func generateHTTPRequestArguments(file *GoFile, args []*types.Argument) {
	for _, arg := range args {
		argName := strings.ToUpperFirst(arg.Name)
		lowerArgName := strings.ToLowerFirst(arg.Name)
		file.Pf("%s %s `json:\"%s\"`", argName, arg.Type.GoType(), lowerArgName)
	}
	file.Pf("")
}

func getArgumentsOfOrigin(args []*types.Argument, origin string) (rArgs []*types.Argument) {
	for _, arg := range args {
		if arg.Options.HTTP.Origin == origin {
			rArgs = append(rArgs, arg)
		}
	}
	return
}
