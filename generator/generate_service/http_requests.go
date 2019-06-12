package generate_service

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateHTTPRequestsFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, method := range helpers.GetMethodsWithHTTPEnabled(service) {
		generateHTTPRequest(file, method)
		generateHTTPRequestNewFunc(file, method)
		generateHTTPRequestNewHTTPFunc(file, method)
		generateHTTPRequestToRequestFunc(file, method)
		generateHTTPRequestToHTTPArgFunc(file, method)
	}
	return file
}

func generateHTTPRequest(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
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

func generateHTTPRequestArguments(file *files.GoFile, args []*types.Argument) {
	for _, arg := range args {
		argName := strings.ToUpperFirst(arg.Name)
		argSpecialName := helpers.GetName(arg.Name, arg.Alias)
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

func generateHTTPRequestNewFunc(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	methodName := strings.ToUpperFirst(method.Name)
	file.AddImport("", method.Service.ImportPath, "service/requests")
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
}

func generateHTTPRequestNewHTTPFunc(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	file.AddImport("", "net/http")
	methodName := strings.ToUpperFirst(method.Name)
	file.Pf("func %sFromHTTP(r *http.Request) (*%sRequest, error) {", methodName, methodName)
	generateHTTPRequestExtractorLogic(file, method)
	file.Pf("}")
	file.Pf("")
}

func generateHTTPRequestToRequestFunc(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	file.AddImport("", method.Service.ImportPath, "service/requests")
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
}

func generateHTTPRequestToHTTPArgFunc(file *files.GoFile, method *types.Method) {
	if len(method.Arguments) == 0 {
		return
	}
	file.AddImport("", "net/http")
	file.AddImport("", method.Service.ImportPath, "service/requests")
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
				argSpecialName := helpers.GetName(arg.Name, arg.Alias)
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
				argSpecialName := strings.ToKebabCase(helpers.GetName(arg.Name, arg.Alias))
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
			file.Pf("goms_http.FormatURI(\"%s\",", getMethodURI(method))
			for _, arg := range getArgumentsOfOrigin(method.Arguments, "PATH") {
				argSpecialName := helpers.GetName(arg.Name, arg.Alias)
				lowerArgName := strings.ToLowerFirst(arg.Name)
				file.Pf("\"%s\", %s,", argSpecialName, lowerArgName)
			}
			file.Pf(")")
		}
	}
	file.Pf("return nil")
	file.Pf("}")
	file.Pf("")
}

func generateHTTPRequestExtractorLogic(file *files.GoFile, method *types.Method) {
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
		generateHTTPRequestExtractors(file, method.Arguments)
	}
	file.Pf("return req, err")
}

func generateHTTPRequestExtractors(file *files.GoFile, args []*types.Argument) {
	if hasArgumentsOfOrigin(args, "QUERY") {
		file.Pf("query := r.URL.Query()")
		for _, arg := range getArgumentsOfOrigin(args, "QUERY") {
			generateQueryExtractor(file, arg)
		}
	}
	if hasArgumentsOfOrigin(args, "HEADER") {
		file.Pf("header := r.Header")
		for _, arg := range getArgumentsOfOrigin(args, "HEADER") {
			generateHeaderExtractor(file, arg)
		}
	}
	if hasArgumentsOfOrigin(args, "PATH") {
		file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
		file.Pf("pathParams:=goms_http.GetParams(r.Context())")
		for _, arg := range getArgumentsOfOrigin(args, "PATH") {
			generatePathExtractor(file, arg)
		}
	}
}

func generateQueryExtractor(file *files.GoFile, arg *types.Argument) {
	argName := strings.ToUpperFirst(arg.Name)
	lowerArgName := strings.ToLowerFirst(arg.Name)
	argNameSnake := strings.ToSnakeCase(helpers.GetName(arg.Name, arg.Alias))
	if arg.Type.IsVariadic || arg.Type.IsSlice {
		file.Pf("%ss := query[\"%s\"]", lowerArgName, argNameSnake)
		file.Pf("if len(%ss)>0 {", lowerArgName)
		file.Pf("var values %s", arg.Type.GoType())
		file.Pf("for _,v:=range %ss {", lowerArgName)
		file.Pf("if len(v) > 0 {")
		generateArgumentConverter(file, arg, "vs", "v")
		file.Pf("values = append(values, %s)", getTypeConverterWrapper(arg.Type.GoType(), "vs"))
		file.Pf("}")
		file.Pf("}")
		file.Pf("req.%s = values", argName)
		file.Pf("}")
	} else {
		file.Pf("v:=query.Get(\"%s\")", argNameSnake)
		file.Pf("if len(v) > 0 {")
		generateArgumentConverter(file, arg, "vs", "v")
		file.Pf("req.%s = vs", argName)
		file.Pf("}")
	}
}

func generateHeaderExtractor(file *files.GoFile, arg *types.Argument) {
	argName := strings.ToUpperFirst(arg.Name)
	lowerArgName := strings.ToLowerFirst(arg.Name)
	argNameKebabCase := strings.ToKebabCase(helpers.GetName(arg.Name, arg.Alias))
	if arg.Type.IsVariadic || arg.Type.IsSlice {
		file.Pf("%ss := header[\"%s\"]", lowerArgName, argNameKebabCase)
		file.Pf("if len(%ss)>0 {", lowerArgName)
		file.Pf("var values %s", arg.Type.GoType())
		file.Pf("for _,v:=range %ss {", lowerArgName)
		file.Pf("if len(v) > 0 {")
		generateArgumentConverter(file, arg, "vs", "v")
		file.Pf("values = append(values, %s)", getTypeConverterWrapper(arg.Type.GoType(), "vs"))
		file.Pf("}")
		file.Pf("}")
		file.Pf("req.%s = values", argName)
		file.Pf("}")
	} else {
		file.Pf("v:=header.Get(\"%s\")", argNameKebabCase)
		file.Pf("if len(v) > 0 {")
		generateArgumentConverter(file, arg, "vs", "v")
		file.Pf("req.%s = vs", argName)
		file.Pf("}")
	}
}

func generatePathExtractor(file *files.GoFile, arg *types.Argument) {
	file.Pf("v:=pathParams.Get(\"%s\")", strings.ToSnakeCase(arg.Name))
	file.Pf("if len(v) > 0 {")
	generateArgumentConverter(file, arg, "vs", "v")
	file.Pf("req.%s = vs", strings.ToUpperFirst(arg.Name))
	file.Pf("}")
}

func generateArgumentConverter(file *files.GoFile, arg *types.Argument, varName string, argName string) {
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
