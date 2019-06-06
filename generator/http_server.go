package generator

import (
	"path"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func generateHTTPServerFile(base string, path string, name string, service *types.Service) *GoFile {
	file := NewGoFile(base, path, name, true, false)
	generateHTTPTransportServerInitFunc(file, service)
	return file
}

func generateHTTPTransportServerInitFunc(file *GoFile, service *types.Service) {
	serviceName := strings.ToUpperFirst(service.Name)
	serviceNameSnake := strings.ToSnakeCase(service.Name)
	file.AddImport("kit_http", "github.com/go-kit/kit/transport/http")
	file.AddImport("", service.ImportPath, "/service/transport")
	file.AddImport(serviceNameSnake+"_http", service.ImportPath, "/service/transport/http")
	file.AddImport("goms_http", "github.com/wlMalk/goms/goms/http")
	file.Pf("func Init(router goms_http.Router, endpoints *transport.%s, opts ...kit_http.ServerOption) {", serviceName)
	serviceVersion := "v" + service.Version.String()
	serviceHTTPURIPrefix := service.Options.HTTP.URIPrefix
	if serviceHTTPURIPrefix == "" {
		serviceHTTPURIPrefix = serviceNameSnake
	}
	for _, method := range service.Methods {
		methodName := strings.ToUpperFirst(method.Name)
		methodNameSnake := strings.ToSnakeCase(method.Name)
		methodHTTPMethod := method.Options.HTTP.Method
		methodHTTPURI := method.Options.HTTP.URI
		methodHTTPabsURI := method.Options.HTTP.AbsURI
		if methodHTTPabsURI == "" {
			if methodHTTPURI == "" {
				methodHTTPURI = methodNameSnake
			}
			methodHTTPabsURI = path.Join(serviceHTTPURIPrefix, methodHTTPURI)
		}
		file.Pf("router.Method(\"%s\", \"/%s\", kit_http.NewServer(", methodHTTPMethod, path.Join(serviceVersion, methodHTTPabsURI))
		file.Pf("endpoints.%s,", methodName)
		file.Pf("%s_http.Decode%sRequest,", serviceNameSnake, methodName)
		file.Pf("%s_http.Encode%sResponse,", serviceNameSnake, methodName)
		file.Pf("opts...))")
	}
	file.Pf("}")
}
