package generator

import (
	"fmt"
	"path/filepath"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func getMethodArguments(args []*types.Argument) []string {
	var a []string
	for _, arg := range args {
		a = append(a, fmt.Sprintf("%s %s", strings.ToLowerFirst(arg.Name), arg.Type.GoArgumentType()))
	}
	return a
}

func getMethodResults(results []*types.Field) []string {
	var r []string
	for _, result := range results {
		r = append(r, fmt.Sprintf("%s %s", strings.ToLowerFirst(result.Name), result.Type.GoArgumentType()))
	}
	return r
}

func addMethodImports(file *GoFile, method *types.Method) {
	file.AddImport("", "context")
}

func generateFunc(file *GoFile, receiver string, name string, args []string, results []string, f func()) {
	file.P(getFuncSignature(receiver, name, args, results))
	f()
	file.P("}")
	file.P("")
}

func getFuncSignature(receiver string, name string, args []string, results []string) string {
	signature := "func "
	if receiver != "" {
		signature += "(" + receiver + ") "
	}
	signature += name + "(" + strs.Join(args, ", ") + ") "
	if len(results) > 0 {
		signature += "(" + strs.Join(results, ", ") + ") "
	}
	signature += "{"
	return signature
}

func getExportedMethodSignature(receiver string, method *types.Method) string {
	return getFuncSignature(receiver, strings.ToUpperFirst(method.Name), getMethodArguments(method.Arguments), getMethodResults(method.Results))
}

func getUnexportedMethodSignature(receiver string, method *types.Method) string {
	return getFuncSignature(receiver, strings.ToLowerFirst(method.Name), getMethodArguments(method.Arguments), getMethodResults(method.Results))
}

func getFieldTagsString(tags map[string]string) string {
	var a []string
	for k, v := range tags {
		a = append(a, fmt.Sprintf("%s:\"%s\"", k, v))
	}
	return strs.Join(a, ",")
}

func generateStruct(file *GoFile, name string, fields []*types.Field) {
	file.Pf("type %s struct {", name)
	for _, f := range fields {
		jsonName := f.Name
		if f.Alias != "" {
			jsonName = f.Alias
		}
		if len(f.Tags) == 0 {
			f.Tags = map[string]string{"json": strings.ToLowerFirst(jsonName)}
		} else {
			f.Tags["json"] = strings.ToLowerFirst(jsonName)
		}
		file.Pf("%s %s `%s`", f.Name, f.Type.GoType(), getFieldTagsString(f.Tags))
	}
	file.P("}")
	file.P("")
}

func generateExportedStruct(file *GoFile, name string, fields []*types.Field) {
	generateStruct(file, strings.ToUpperFirst(name), fields)
}

func generateUnexportedStruct(file *GoFile, name string, fields []*types.Field) {
	generateStruct(file, strings.ToLowerFirst(name), fields)
}

func getMethodArgumentsInCall(args []*types.Argument) (a []string) {
	for _, arg := range args {
		if arg.Type.IsVariadic {
			a = append(a, strings.ToLowerFirst(arg.Name)+"...")
		} else {
			a = append(a, strings.ToLowerFirst(arg.Name))
		}
	}
	return
}

func getMethodArgumentsFromRequestInCall(args []*types.Argument) (a []string) {
	a = getMethodArgumentsInCall(args)
	for i := range a {
		a[i] = "req." + strings.ToUpperFirst(a[i])
	}
	return
}

func getResultsVars(results []*types.Field) (a []string) {
	for _, result := range results {
		a = append(a, strings.ToLowerFirst(result.Name))
	}
	return
}

func getResultsVarsFromResponse(results []*types.Field) (a []string) {
	a = getResultsVars(results)
	for i := range a {
		a[i] = "res." + strings.ToUpperFirst(a[i])
	}
	return
}

func Generate(s *types.Service) (files Files, err error) {
	files = append(files, generateRequestsFile(s.Path, filepath.Join("service", "requests"), "requests", s.Methods))
	files = append(files, generateResponseFile(s.Path, filepath.Join("service", "responses"), "responses", s.Methods))
	files = append(files, generateHandlersFile(s.Path, filepath.Join("service", "handlers"), "handlers", s))
	files = append(files, generateConvertersFile(s.Path, filepath.Join("service", "handlers", "converters"), "converters", s.Methods))
	files = append(files, generateServiceMiddlewareFile(s.Path, filepath.Join("service", "middleware"), "middleware", s))
	files = append(files, generateServiceStructFile(s.Path, filepath.Join("service"), strings.ToLowerFirst(s.Name), s))
	return
}
