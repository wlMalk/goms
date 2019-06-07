package generator

import (
	"path"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPServerFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateHTTPTransportServerStruct(file, service)
	generateHTTPTransportServerNewFunc(file, service)
	return file
}

func generateHTTPTransportServerStruct(file *GoFile, service *types.Service) {
	file.Pf("type Server struct {")
	file.Pf("http.Server")
	file.Pf("}")
	file.Pf("")
}

func generateHTTPTransportServerNewFunc(file *GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("", "net/http")
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
	file.AddImport("", service.ImportPath, "/service/transport")
	file.AddImport(serviceNameSnake+"_http", service.ImportPath, "/service/transport/http")
	file.AddImport("goms_http", "github.com/wlMalk/goms/goms/http")
	file.Pf("func New(router goms_http.Router, endpoints *transport.%s, opts ...kit_http.ServerOption) *Server {", serviceName)
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		methodHTTPMethod := method.Options.HTTP.Method
		methodURI := getMethodURI(method)
		file.Pf("router.Method(\"%s\", \"%s\", kit_http.NewServer(", methodHTTPMethod, methodURI)
		file.Pf("endpoints.%s,", methodName)
		file.Pf("%s_http.Decode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_http.Encode%sResponse,", serviceNameSnake, methodName)
		file.Pf("opts...))")
	}
	file.Pf("return &Server{Server: http.Server{Handler: router}}")
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
