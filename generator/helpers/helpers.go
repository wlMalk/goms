package helpers

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GetMethodArguments(args []*types.Argument) []string {
	var a []string
	for _, arg := range args {
		a = append(a, fmt.Sprintf("%s %s", strings.ToLowerFirst(arg.Name), arg.Type.GoArgumentType()))
	}
	return a
}

func GetMethodResults(results []*types.Field) []string {
	var r []string
	for _, result := range results {
		r = append(r, fmt.Sprintf("%s %s", strings.ToLowerFirst(result.Name), result.Type.GoArgumentType()))
	}
	return r
}

func AddMethodImports(file *files.GoFile, method *types.Method) {
	file.AddImport("", "context")
}

func GenerateFunc(file *files.GoFile, receiver string, name string, args []string, results []string, f func()) {
	file.P(GetFuncSignature(receiver, name, args, results))
	f()
	file.P("}")
	file.P("")
}

func GetFuncSignature(receiver string, name string, args []string, results []string) string {
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

func GetExportedMethodSignature(receiver string, method *types.Method) string {
	return GetFuncSignature(receiver, strings.ToUpperFirst(method.Name), GetMethodArguments(method.Arguments), GetMethodResults(method.Results))
}

func GetUnexportedMethodSignature(receiver string, method *types.Method) string {
	return GetFuncSignature(receiver, strings.ToLowerFirst(method.Name), GetMethodArguments(method.Arguments), GetMethodResults(method.Results))
}

func GetFieldTagsString(tags map[string]string) string {
	var a []string
	for k, v := range tags {
		a = append(a, fmt.Sprintf("%s:\"%s\"", k, v))
	}
	return strs.Join(a, ",")
}

func GenerateStruct(file *files.GoFile, name string, fields []*types.Field) {
	if len(fields) == 0 {
		return
	}
	file.Pf("type %s struct {", name)
	for _, f := range fields {
		jsonName := GetName(f.Name, f.Alias)
		if len(f.Tags) == 0 {
			f.Tags = map[string]string{"json": strings.ToLowerFirst(jsonName)}
		} else {
			f.Tags["json"] = strings.ToLowerFirst(jsonName)
		}
		file.Pf("%s %s `%s`", f.Name, f.Type.GoType(), GetFieldTagsString(f.Tags))
	}
	file.P("}")
	file.P("")
}

func GenerateExportedStruct(file *files.GoFile, name string, fields []*types.Field) {
	GenerateStruct(file, strings.ToUpperFirst(name), fields)
}

func GenerateUnexportedStruct(file *files.GoFile, name string, fields []*types.Field) {
	GenerateStruct(file, strings.ToLowerFirst(name), fields)
}

func GetMethodArgumentsInCall(args []*types.Argument) (a []string) {
	for _, arg := range args {
		if arg.Type.IsVariadic {
			a = append(a, strings.ToLowerFirst(arg.Name)+"...")
		} else {
			a = append(a, strings.ToLowerFirst(arg.Name))
		}
	}
	return
}

func GetMethodArgumentsFromRequestInCall(args []*types.Argument) (a []string) {
	a = GetMethodArgumentsInCall(args)
	for i := range a {
		a[i] = "req." + strings.ToUpperFirst(a[i])
	}
	return
}

func GetResultsVars(results []*types.Field) (a []string) {
	for _, result := range results {
		a = append(a, strings.ToLowerFirst(result.Name))
	}
	return
}

func GetResultsVarsFromResponse(results []*types.Field) (a []string) {
	a = GetResultsVars(results)
	for i := range a {
		a[i] = "res." + strings.ToUpperFirst(a[i])
	}
	return
}

func GetName(name string, alias string) string {
	n := name
	if alias != "" {
		n = alias
	}
	return strings.ToLowerFirst(n)
}
