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

func ParseVersion(ver string) (*types.Version, error) {
	v := &types.Version{}
	var err error
	v.Major, v.Minor, v.Patch, err = parseVersion(ver)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func Parse(ast *astTypes.File) (s *types.Service, err error) {
	s = &types.Service{}
	s.Name = strings.ToUpperFirst(strings.ToCamelCase(ast.Name))
	iface := getServiceInterface(ast.Interfaces)
	if iface == nil {
		return nil, errors.New("Service interface is not defined in service.go file")
	}
	for _, method := range iface.Methods {
		m, err := parseMethod(method)
		if err != nil {
			return nil, err
		}
		if err := validateMethod(m); err != nil {
			return nil, err
		}
		s.Methods = append(s.Methods, m)
	}
	s.Tags, s.Docs = cleanComments(iface.Docs)
	// ast.
	return
}

func parseArgument(v astTypes.Variable) (*types.Argument, error) {
	a := &types.Argument{}
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
			if err := validateField(r); err != nil {
				return err
			}
			m.Results = append(m.Results, r)
		}
	}
	return nil
}

func parseMethod(method *astTypes.Function) (*types.Method, error) {
	m := &types.Method{}
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
