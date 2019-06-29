package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func HTTPRequest(file file.File, service types.Service, method types.Method) error {
	helpers.AddTypesImports(file, service)
	if len(method.Arguments) == 0 {
		return nil
	}
	if (method.Options.HTTP.Method == "POST" || method.Options.HTTP.Method == "PUT") && hasArgumentsOfOrigin(method.Arguments, "BODY") {
		file.Pf("type " + method.Name + "RequestBody struct {")
		HTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "BODY"))
		file.Pf("}")
		file.Pf("")
	}
	file.Pf("type " + method.Name + "Request struct {")
	if (method.Options.HTTP.Method == "POST" || method.Options.HTTP.Method == "PUT") && hasArgumentsOfOrigin(method.Arguments, "BODY") {
		file.Pf("Body *" + method.Name + "RequestBody")
		file.Pf("")
	}
	if hasArgumentsOfOrigin(method.Arguments, "HEADER") {
		HTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "HEADER"))
	}
	if hasArgumentsOfOrigin(method.Arguments, "QUERY") {
		HTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "QUERY"))
	}
	if hasArgumentsOfOrigin(method.Arguments, "PATH") {
		HTTPRequestArguments(file, getArgumentsOfOrigin(method.Arguments, "PATH"))
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func hasArgumentsOfOrigin(args []*types.Argument, origin string) bool {
	for _, arg := range args {
		if arg.Options.HTTP.Origin == origin {
			return true
		}
	}
	return false
}

func HTTPRequestArguments(file file.File, args []*types.Argument) {
	for _, arg := range args {
		argName := strings.ToUpperFirst(arg.Name)
		argSpecialName := helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias)
		file.Pf("%s %s `json:\"%s\"`", argName, arg.Type.GoType(), argSpecialName)
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

func HTTPRequestNewFunc(file file.File, service types.Service, method types.Method) error {
	if len(method.Arguments) == 0 {
		return nil
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	file.Pf("func %s(req *requests.%sRequest) *%sRequest {", methodName, methodName, methodName)
	file.Pf("r := &%sRequest{}", methodName)
	for _, arg := range method.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		if arg.Options.HTTP.Origin == "BODY" {
			file.Pf("r.Body.%s = req.%s", argName, argName)
		} else {
			file.Pf("r.%s = req.%s", argName, argName)
		}
	}
	file.Pf("return r")
	file.Pf("}")
	file.Pf("")
	return nil
}

func HTTPRequestNewHTTPFunc(file file.File, service types.Service, method types.Method) error {
	if len(method.Arguments) == 0 {
		return nil
	}
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func %sFromHTTP(r *http.Request) (*%sRequest, error) {", methodName, methodName)
	HTTPRequestExtractorLogic(file, service, method)
	file.Pf("}")
	file.Pf("")
	return nil
}

func HTTPRequestToRequestFunc(file file.File, service types.Service, method types.Method) error {
	if len(method.Arguments) == 0 {
		return nil
	}
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func (r *%sRequest) Request() *requests.%sRequest {", methodName, methodName)
	file.Pf("req := &requests.%sRequest{}", methodName)
	for _, arg := range method.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		if arg.Options.HTTP.Origin == "BODY" {
			file.Pf("req.%s = r.Body.%s", argName, argName)
		} else {
			file.Pf("req.%s = r.%s", argName, argName)
		}
	}
	file.Pf("return req")
	file.Pf("}")
	file.Pf("")
	return nil
}

func HTTPRequestToHTTPArgFunc(file file.File, service types.Service, method types.Method) error {
	if len(method.Arguments) == 0 {
		return nil
	}
	file.AddImport("", "net/http")
	file.AddImport("", service.ImportPath, "/pkg/service/requests")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func (r *%sRequest) ToHTTP(req *http.Request) error {", methodName)
	file.Pf("var err error")
	if (method.Options.HTTP.Method == "POST" || method.Options.HTTP.Method == "PUT") && hasArgumentsOfOrigin(method.Arguments, "BODY") {
		file.AddImport("", "encoding/json")
		file.AddImport("", "io/ioutil")
		file.AddImport("", "bytes")
		file.Pf("var buf bytes.Buffer")
		file.Pf("if err = json.NewEncoder(&buf).Encode(r.Body); err != nil {")
		file.Pf("return err")
		file.Pf("}")
		file.Pf("req.Body = ioutil.NopCloser(&buf)")
	}
	if hasArgumentsOfOrigin(method.Arguments, "HEADER") ||
		hasArgumentsOfOrigin(method.Arguments, "QUERY") ||
		hasArgumentsOfOrigin(method.Arguments, "PATH") {
		file.AddImport("goms_util", "github.com/wlMalk/goms/goms/util")
		if hasArgumentsOfOrigin(method.Arguments, "HEADER") ||
			hasArgumentsOfOrigin(method.Arguments, "QUERY") {
			file.Pf("var value string")
		}
		if hasArgumentsOfOrigin(method.Arguments, "QUERY") {
			file.Pf("query := req.URL.Query()")
			for _, arg := range getArgumentsOfOrigin(method.Arguments, "QUERY") {
				argName := strings.ToUpperFirst(arg.Name)
				argSpecialName := helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias)
				if arg.Type.IsSlice || arg.Type.IsVariadic {
					file.Pf("for i := range r.%s {", argName)
					file.Pf("value, err = goms_util.ToString(r.%s[i])", argName)
					file.Pf("if err != nil {")
					file.Pf("return err")
					file.Pf("}")
					file.Pf("query.Add(\"%s\", value)", argSpecialName)
					file.Pf("}")
				} else {
					file.Pf("value, err = goms_util.ToString(r.%s)", argName)
					file.Pf("if err != nil {")
					file.Pf("return err")
					file.Pf("}")
					file.Pf("query.Add(\"%s\", value)", argSpecialName)
				}
			}
		}
		if hasArgumentsOfOrigin(method.Arguments, "HEADER") {
			file.Pf("header := req.Header")
			for _, arg := range getArgumentsOfOrigin(method.Arguments, "HEADER") {
				argName := strings.ToUpperFirst(arg.Name)
				argSpecialName := strings.ToKebabCase(helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias))
				if arg.Type.IsSlice || arg.Type.IsVariadic {
					file.Pf("for i := range r.%s {", argName)
					file.Pf("value, err = goms_util.ToString(r.%s[i])", argName)
					file.Pf("if err != nil {")
					file.Pf("return err")
					file.Pf("}")
					file.Pf("header.Add(\"%s\", value)", argSpecialName)
					file.Pf("}")
				} else {
					file.Pf("value, err = goms_util.ToString(r.%s)", argName)
					file.Pf("if err != nil {")
					file.Pf("return err")
					file.Pf("}")
					file.Pf("header.Add(\"%s\", value)", argSpecialName)
				}
			}
		}
		if hasArgumentsOfOrigin(method.Arguments, "PATH") {
			file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
			for _, arg := range method.Arguments {
				argName := strings.ToUpperFirst(arg.Name)
				lowerArgName := strings.ToLowerFirst(arg.Name)
				file.Pf("%s, err := goms_util.ToString(r.%s)", lowerArgName, argName)
				file.Pf("if err != nil {")
				file.Pf("return err")
				file.Pf("}")
			}
			file.Pf("goms_http.FormatURI(\"%s\",", getMethodURI(service, method))
			for _, arg := range getArgumentsOfOrigin(method.Arguments, "PATH") {
				argSpecialName := helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias)
				lowerArgName := strings.ToLowerFirst(arg.Name)
				file.Pf("\"%s\", %s,", argSpecialName, lowerArgName)
			}
			file.Pf(")")
		}
	}
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
	return nil
}

func HTTPRequestExtractorLogic(file file.File, service types.Service, method types.Method) error {
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("var err error")
	file.Pf("req := &%sRequest{}", methodName)
	if (method.Options.HTTP.Method == "POST" || method.Options.HTTP.Method == "PUT") && hasArgumentsOfOrigin(method.Arguments, "BODY") {
		file.AddImport("", "encoding/json")
		file.Pf("d := json.NewDecoder(r.Body)")
		file.Pf("err = d.Decode(&req.Body)")
		file.Pf("if err != nil {")
		file.Pf("return nil, err")
		file.Pf("}")
	}
	if hasArgumentsOfOrigin(method.Arguments, "HEADER") ||
		hasArgumentsOfOrigin(method.Arguments, "QUERY") ||
		hasArgumentsOfOrigin(method.Arguments, "PATH") {
		HTTPRequestExtractors(file, method.Arguments)
	}
	file.Pf("return req, err")
	return nil
}

func HTTPRequestExtractors(file file.File, args []*types.Argument) {
	if hasArgumentsOfOrigin(args, "QUERY") {
		file.Pf("query := r.URL.Query()")
		for _, arg := range getArgumentsOfOrigin(args, "QUERY") {
			QueryExtractor(file, arg)
		}
	}
	if hasArgumentsOfOrigin(args, "HEADER") {
		file.Pf("header := r.Header")
		for _, arg := range getArgumentsOfOrigin(args, "HEADER") {
			HeaderExtractor(file, arg)
		}
	}
	if hasArgumentsOfOrigin(args, "PATH") {
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
		file.Pf("pathParams:=goms_http.GetParams(r.Context())")
		for _, arg := range getArgumentsOfOrigin(args, "PATH") {
			PathExtractor(file, arg)
		}
	}
}

func QueryExtractor(file file.File, arg *types.Argument) {
	argName := strings.ToUpperFirst(arg.Name)
	lowerArgName := strings.ToLowerFirst(arg.Name)
	argNameSnake := strings.ToSnakeCase(helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias))
	if arg.Type.IsVariadic || arg.Type.IsSlice {
		file.Pf("%ss := query[\"%s\"]", lowerArgName, argNameSnake)
		file.Pf("if len(%ss)>0 {", lowerArgName)
		file.Pf("var values %s", arg.Type.GoType())
		file.Pf("for _,v:=range %ss {", lowerArgName)
		file.Pf("if len(v) > 0 {")
		ArgumentConverter(file, arg, "vs", "v")
		file.Pf("values = append(values, %s)", getTypeConverterWrapper(arg.Type.GoType(), "vs"))
		file.Pf("}")
		file.Pf("}")
		file.Pf("req.%s = values", argName)
		file.Pf("}")
	} else {
		file.Pf("v:=query.Get(\"%s\")", argNameSnake)
		file.Pf("if len(v) > 0 {")
		ArgumentConverter(file, arg, "vs", "v")
		file.Pf("req.%s = vs", argName)
		file.Pf("}")
	}
}

func HeaderExtractor(file file.File, arg *types.Argument) {
	argName := strings.ToUpperFirst(arg.Name)
	lowerArgName := strings.ToLowerFirst(arg.Name)
	argNameKebabCase := strings.ToKebabCase(helpers.GetName(strings.ToLowerFirst(arg.Name), arg.Alias))
	if arg.Type.IsVariadic || arg.Type.IsSlice {
		file.Pf("%ss := header[\"%s\"]", lowerArgName, argNameKebabCase)
		file.Pf("if len(%ss)>0 {", lowerArgName)
		file.Pf("var values %s", arg.Type.GoType())
		file.Pf("for _,v:=range %ss {", lowerArgName)
		file.Pf("if len(v) > 0 {")
		ArgumentConverter(file, arg, "vs", "v")
		file.Pf("values = append(values, %s)", getTypeConverterWrapper(arg.Type.GoType(), "vs"))
		file.Pf("}")
		file.Pf("}")
		file.Pf("req.%s = values", argName)
		file.Pf("}")
	} else {
		file.Pf("v:=header.Get(\"%s\")", argNameKebabCase)
		file.Pf("if len(v) > 0 {")
		ArgumentConverter(file, arg, "vs", "v")
		file.Pf("req.%s = vs", argName)
		file.Pf("}")
	}
}

func PathExtractor(file file.File, arg *types.Argument) {
	file.Pf("v:=pathParams.Get(\"%s\")", strings.ToSnakeCase(arg.Name))
	file.Pf("if len(v) > 0 {")
	ArgumentConverter(file, arg, "vs", "v")
	file.Pf("req.%s = vs", strings.ToUpperFirst(arg.Name))
	file.Pf("}")
}

func ArgumentConverter(file file.File, arg *types.Argument, varName string, argName string) {
	argNameSnake := strings.ToSnakeCase(arg.Name)
	switch getBasicType(arg.Type.Name) {
	case "int":
		file.Pf("%s,err:=strconv.ParseInt(%s, 10, 64)", varName, argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidNumber(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
	case "float":
		file.Pf("%s,err:=strconv.ParseFloat(%s, 10, 64)", varName, argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidNumber(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
	case "bool":
		file.Pf("%s,err:=strconv.ParseBool(%s, 10, 64)", varName, argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidBool(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
	case "uint":
		file.Pf("%s,err:=strconv.ParseUint(%s, 10, 64)", varName, argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidNumber(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
	case "time":
		file.Pf("timestamp,err:=strconv.ParseInt(%s, 10, 64)", argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidUnixTime(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
		file.Pf("%s:=time.Unix(timestamp, 0)", varName)
	case "uuid":
		file.Pf("%s, err:=uuid.FromString(%s)", varName, argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidUUID(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
	case "duration":
		file.Pf("duration,err:=strconv.ParseInt(%s, 10, 64)", argName)
		file.Pf("if err!=nil{")
		file.Pf("return errors.ErrParamNotValidDuration(\"%s\",%s)", argNameSnake, argName)
		file.Pf("}")
		file.Pf("%s := time.Duration(duration)", varName)
	default:
		file.Pf("%s := %s", varName, argName)
	}
}

func getTypeConverterWrapper(t string, n string) string {
	if strs.Contains(t, "int") || strs.Contains(t, "float") || strs.Contains(t, "uint") {
		if t == "int64" || t == "float64" || t == "uint64" {
			return n
		}
		return t + "(" + n + ")"
	}
	return n
}

func getBasicType(t string) string {
	if strs.Contains(t, "uint") {
		return "uint"
	}
	if strs.Contains(t, "float") {
		return "float"
	}
	if strs.Contains(t, "int") {
		return "int"
	}
	return strs.Replace(t, "[]", "", -1)
}
