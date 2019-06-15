package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"

	astTypes "github.com/vetcher/go-astra/types"
)

var versionPattern = regexp.MustCompile(`(?is)^v?([0-9]+)((\.|-|_|:)([0-9]+))?((\.|-|_|:)([0-9]+))?$`)

func Parse(ast *astTypes.File) (services []*types.Service, err error) {
	serviceName := cleanServiceName(ast.Name)
	ifaces := getServiceInterfaces(ast.Interfaces, serviceName)
	if len(ifaces) == 0 {
		return nil, errors.New("no service definitions were found")
	}
	for _, iface := range ifaces {
		s := &types.Service{}
		s.Name = serviceName
		ver, err := ParseVersion(extractServiceVersion(iface.Name, serviceName))
		if err != nil {
			return nil, err
		}
		s.Version = *ver
		var tags []string
		tags, s.Docs = cleanComments(iface.Docs)
		err = parseServiceTags(s, tags)
		if err != nil {
			return nil, err
		}
		for _, method := range iface.Methods {
			m, tags, err := parseMethod(method)
			if err != nil {
				return nil, err
			}
			m.Service = s
			m.Options.Generate.Caching = s.Options.Generate.Caching
			m.Options.Generate.Logging = s.Options.Generate.Logging
			m.Options.Generate.MethodStubs = s.Options.Generate.MethodStubs
			m.Options.Generate.Middleware = s.Options.Generate.Middleware
			m.Options.Generate.Validators = s.Options.Generate.Validators
			m.Options.Generate.Validating = s.Options.Generate.Validating
			m.Options.Generate.CircuitBreaking = s.Options.Generate.CircuitBreaking
			m.Options.Generate.RateLimiting = s.Options.Generate.RateLimiting
			m.Options.Generate.Recovering = s.Options.Generate.Recovering
			m.Options.Generate.Tracing = s.Options.Generate.Tracing
			m.Options.Generate.FrequencyMetric = s.Options.Generate.FrequencyMetric
			m.Options.Generate.LatencyMetric = s.Options.Generate.LatencyMetric
			m.Options.Generate.CounterMetric = s.Options.Generate.CounterMetric
			m.Options.Generate.HTTPServer = s.Options.Generate.HTTPServer
			m.Options.Generate.HTTPClient = s.Options.Generate.HTTPClient
			m.Options.Generate.GRPCServer = s.Options.Generate.GRPCServer
			m.Options.Generate.GRPCClient = s.Options.Generate.GRPCClient
			err = parseMethodTags(m, tags)
			if err != nil {
				return nil, err
			}
			if err := validateMethod(m); err != nil {
				return nil, err
			}
			s.Methods = append(s.Methods, m)
		}
		services = append(services, s)
	}
	return
}

func ParseVersion(ver string) (*types.Version, error) {
	v := &types.Version{}
	var err error
	v.Major, v.Minor, v.Patch, err = parseVersion(strs.ToLower(ver))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func parseVersion(ver string) (int, int, int, error) {
	err := fmt.Errorf("cannot parse \"%s\" as a version", ver)
	if ver == "" {
		return 1, 0, 0, nil
	}
	matches := versionPattern.FindAllStringSubmatch(ver, -1)
	if len(matches) != 1 {
		return 0, 0, 0, err
	}
	major := 0
	minor := 0
	patch := 0
	var nerr error
	if matches[0][1] != "" {
		if major, nerr = strconv.Atoi(matches[0][1]); nerr != nil {
			return 0, 0, 0, err
		}
	}
	if matches[0][4] != "" {
		if minor, nerr = strconv.Atoi(matches[0][4]); nerr != nil {
			return 0, 0, 0, err
		}
	}
	if matches[0][7] != "" {
		if patch, nerr = strconv.Atoi(matches[0][7]); nerr != nil {
			return 0, 0, 0, err
		}
	}
	return major, minor, patch, nil
}

func defaultArgument() *types.Argument {
	a := &types.Argument{}
	a.Options.HTTP.Origin = "BODY"
	return a
}

func parseArgument(v astTypes.Variable) (*types.Argument, error) {
	a := defaultArgument()
	var err error
	a.Name = strings.ToUpperFirst(strings.ToCamelCase(v.Name))
	a.Type, err = parseType(v.Type)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func parseField(v astTypes.Variable) (*types.Field, error) {
	f := &types.Field{}
	var err error
	f.Name = strings.ToUpperFirst(strings.ToCamelCase(v.Name))
	f.Type, err = parseType(v.Type)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func parseArguments(m *types.Method, args []astTypes.Variable) error {
	for i, arg := range args {
		if i == 0 {
			firstIsContext := true
			t, ok := arg.Type.(astTypes.TImport)
			if !ok || t.Import.Package != "context" || t.Import.Name != "context" {
				firstIsContext = false
			}
			tt, ok := t.Next.(astTypes.TName)
			if !ok || tt.TypeName != "Context" {
				firstIsContext = false
			}
			if !firstIsContext {
				return fmt.Errorf("first argument in \"%s\" method has to be of type \"context.Context\" from package \"context\"", m.Name)
			}
		} else {
			a, err := parseArgument(arg)
			if err != nil {
				return err
			}
			if a.Name == "" {
				return fmt.Errorf("'%s' method has an unnamed argument", m.Name)
			}
			if err := validateArgument(a); err != nil {
				return err
			}
			m.Arguments = append(m.Arguments, a)
		}
	}
	return nil
}

func parseResults(m *types.Method, args []astTypes.Variable) error {
	for i, arg := range args {
		if i == len(args)-1 {
			t, ok := arg.Type.(astTypes.TName)
			if !ok || t.TypeName != "error" {
				return fmt.Errorf("last return value in \"%s\" method has to be of type \"error\"", m.Name)
			}
		} else {
			r, err := parseField(arg)
			if err != nil {
				return err
			}
			if r.Name == "" {
				return fmt.Errorf("'%s' method has an unnamed return value", m.Name)
			}
			if err := validateField(r); err != nil {
				return err
			}
			m.Results = append(m.Results, r)
		}
	}
	return nil
}

func defaultMethod() *types.Method {
	m := &types.Method{}
	m.Options.HTTP.Method = "POST"
	return m
}

func parseMethod(method *astTypes.Function) (*types.Method, []string, error) {
	m := defaultMethod()
	m.Name = method.Name
	err := parseArguments(m, method.Args)
	if err != nil {
		return nil, nil, err
	}
	err = parseResults(m, method.Results)
	if err != nil {
		return nil, nil, err
	}
	var tags []string
	tags, m.Docs = cleanComments(method.Docs)
	return m, tags, nil
}

func validateMethod(m *types.Method) error {
	return nil
}

func validateArgument(a *types.Argument) error {
	return nil
}

func validateField(f *types.Field) error {
	return nil
}

func getServiceInterfaces(interfaces []astTypes.Interface, serviceName string) (ifaces []*astTypes.Interface) {
	for _, i := range interfaces {
		if isServiceInterface(i.Name, serviceName) {
			ifaces = append(ifaces, &i)
		}
	}
	return
}

func isServiceInterface(name string, serviceName string) bool {
	name = strings.ToUpperFirst(name)
	return strs.HasPrefix(name, serviceName+"Service") ||
		strs.HasPrefix(name, "Service") ||
		strs.HasPrefix(name, serviceName) ||
		strs.HasSuffix(name, serviceName+"Service") ||
		strs.HasSuffix(name, "Service") ||
		strs.HasSuffix(name, serviceName) ||
		strs.HasPrefix(name, serviceName+"Svc") ||
		strs.HasPrefix(name, "Svc") ||
		strs.HasPrefix(name, serviceName) ||
		strs.HasSuffix(name, serviceName+"Svc") ||
		strs.HasSuffix(name, "Svc") ||
		strs.HasSuffix(name, serviceName)
}

func extractServiceVersion(name string, serviceName string) string {
	name = strings.ToUpperFirst(name)
	name = strs.Replace(name, serviceName, "", 1)
	name = strs.Replace(name, "Service", "", 1)
	name = strs.Replace(name, "Svc", "", 1)
	name = strs.TrimLeft(name, "_")
	name = strs.TrimRight(name, "_")
	return name
}

func cleanServiceName(name string) string {
	name = strs.TrimSuffix(name, "_service")
	name = strs.TrimSuffix(name, "_svc")
	name = strings.ToCamelCase(name)
	name = strings.ToUpperFirst(name)
	return name
}
