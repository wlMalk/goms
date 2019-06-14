package generate_service

import (
	"path"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateHTTPServerFile(base string, path string, name string, service *types.Service) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	generateHTTPTransportServerRegisterFunc(file, service)
	generateHTTPTransportServerRegisterSpecialFunc(file, service)
	return file
}

func generateHTTPTransportServerRegisterFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
	file.AddImport("", service.ImportPath, "/pkg/service/transport")
	file.AddImport(serviceNameSnake+"_http", service.ImportPath, "/pkg/service/transport/http")
	file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
	file.Pf("func Register(server *goms_http.Server, endpoints *transport.%s, opts ...kit_http.ServerOption) {", serviceName)
	for _, method := range helpers.GetMethodsWithHTTPServerEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		methodHTTPMethod := method.Options.HTTP.Method
		methodURI := getMethodURI(method)
		file.Pf("server.RegisterMethod(\"%s\", \"%s\", kit_http.NewServer(", methodHTTPMethod, methodURI)
		file.Pf("endpoints.%s,", methodName)
		file.Pf("%s_http.Decode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_http.Encode%sResponse,", serviceNameSnake, methodName)
		file.Pf("opts...),")
		file.Pf(")")
	}
	file.Pf("}")
	file.Pf("")
}

func generateHTTPTransportServerRegisterSpecialFunc(file *files.GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
	file.AddImport("", service.ImportPath, "/pkg/service/transport")
	file.AddImport(serviceNameSnake+"_http", service.ImportPath, "/pkg/service/transport/http")
	file.AddImport("goms_http", "github.com/wlMalk/goms/goms/transport/http")
	file.Pf("func RegisterSpecial(server *goms_http.Server, endpoints *transport.%s, optionsFunc func(method string) (opts []kit_http.ServerOption)) {", serviceName)
	for _, method := range helpers.GetMethodsWithHTTPServerEnabled(service) {
		methodName := strings.ToUpperFirst(method.Name)
		methodHTTPMethod := method.Options.HTTP.Method
		methodURI := getMethodURI(method)
		file.Pf("server.RegisterMethod(\"%s\", \"%s\", kit_http.NewServer(", methodHTTPMethod, methodURI)
		file.Pf("endpoints.%s,", methodName)
		file.Pf("%s_http.Decode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_http.Encode%sResponse,", serviceNameSnake, methodName)
		file.Pf("optionsFunc(\"%s\")...),", helpers.GetName(methodName, method.Alias))
		file.Pf(")")
	}
	file.Pf("}")
	file.Pf("")
}

func getMethodURI(method *types.Method) string {
	serviceNameSnake := strings.ToSnakeCase(method.Service.Name)
	serviceVersion := "v" + method.Service.Version.String()
	serviceHTTPURIPrefix := method.Service.Options.HTTP.URIPrefix
	if serviceHTTPURIPrefix == "" {
		serviceHTTPURIPrefix = serviceNameSnake
	}
	methodNameSnake := strings.ToSnakeCase(method.Name)
	methodHTTPURI := method.Options.HTTP.URI
	methodHTTPabsURI := method.Options.HTTP.AbsURI
	if methodHTTPabsURI == "" {
		if methodHTTPURI == "" {
			methodHTTPURI = methodNameSnake
		}
		methodHTTPabsURI = path.Join(serviceHTTPURIPrefix, methodHTTPURI)
	}
	return "/" + path.Join(serviceVersion, methodHTTPabsURI)
}
