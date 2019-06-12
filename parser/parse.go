package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"

	astTypes "github.com/vetcher/go-astra/types"
)

var versionPattern = regexp.MustCompile(`(?is)^v?([0-9]+)((\.|-|_|:)([0-9]+))?((\.|-|_|:)([0-9]+))?$`)

func Parse(ast *astTypes.File) (s *types.Service, err error) {
	s = &types.Service{}
	s.Name = strings.ToUpperFirst(strings.ToCamelCase(ast.Name))
	iface := getServiceInterface(ast.Interfaces)
	if iface == nil {
		return nil, errors.New("Service interface is not defined in service.go file")
	}
	s.Tags, s.Docs = cleanComments(iface.Docs)
	err = parseServiceTags(s, s.Tags)
	if err != nil {
		return nil, err
	}
	for _, method := range iface.Methods {
		m, err := parseMethod(method)
		if err != nil {
			return nil, err
		}
		m.Service = s
		m.Options.Generate.Caching = s.Options.Generate.Caching
		m.Options.Generate.Logging = s.Options.Generate.Logging
		m.Options.Generate.MethodStubs = s.Options.Generate.MethodStubs
		m.Options.Generate.Middleware = s.Options.Generate.Middleware
		m.Options.Generate.Validator = s.Options.Generate.Validators
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
		m.Options.Generate.GRPCServer = s.Options.Generate.GRPCServer
		err = parseMethodTags(m, m.Tags)
		if err != nil {
			return nil, err
		}
		if err := validateMethod(m); err != nil {
			return nil, err
		}
		s.Methods = append(s.Methods, m)
	}
	return
}

func ParseVersion(ver string) (*types.Version, error) {
	v := &types.Version{}
	var err error
	v.Major, v.Minor, v.Patch, err = parseVersion(ver)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func parseVersion(ver string) (int, int, int, error) {
	err := fmt.Errorf("cannot parse \"%s\" as a version", ver)
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

func parseMethod(method *astTypes.Function) (*types.Method, error) {
	m := defaultMethod()
	m.Name = method.Name
	err := parseArguments(m, method.Args)
	if err != nil {
		return nil, err
	}
	err = parseResults(m, method.Results)
	if err != nil {
		return nil, err
	}
	m.Tags, m.Docs = cleanComments(method.Docs)
	return m, nil
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

func getServiceInterface(interfaces []astTypes.Interface) *astTypes.Interface {
	for _, i := range interfaces {
		if i.Name == "Service" {
			return &i
		}
	}
	return nil
}
