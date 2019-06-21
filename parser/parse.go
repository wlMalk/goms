package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"

	astra "github.com/vetcher/go-astra"
	astTypes "github.com/vetcher/go-astra/types"
)

var versionPattern = regexp.MustCompile(`(?is)^v?([0-9]+)((\.|-|_|:)([0-9]+))?((\.|-|_|:)([0-9]+))?$`)

func (p *Parser) Parse(f *ast.File) (services []types.Service, err error) {
	file, err := astra.ParseAstFile(f)
	if err != nil {
		return nil, err
	}
	enums, err := p.parseEnums(f)
	if err != nil {
		return nil, err
	}
	entities, err := p.parseEntities(file)
	if err != nil {
		return nil, err
	}
	serviceName := cleanServiceName(file.Name)
	ifaces := getServiceInterfaces(file.Interfaces, serviceName)
	if len(ifaces) == 0 {
		return nil, errors.New("no service definitions were found")
	}
	for _, iface := range ifaces {
		s, err := p.parseService(iface, serviceName)
		if err != nil {
			return nil, err
		}
		p.setServiceEnums(s, enums)
		p.setServiceEntities(s, entities)
		services = append(services, *s)
	}
	return
}

func (p *Parser) parseService(iface *astTypes.Interface, serviceName string) (*types.Service, error) {
	s := defaultService()
	s.Name = serviceName
	ver, err := ParseVersion(extractServiceVersion(iface.Name, serviceName))
	if err != nil {
		return nil, err
	}
	s.Version = *ver
	var ts []string
	ts, s.Docs = cleanComments(iface.Docs)
	err = p.parseServiceTags(s, ts)
	if err != nil {
		return nil, err
	}
	for _, method := range iface.Methods {
		m, ts, err := p.parseMethod(method)
		if err != nil {
			return nil, err
		}
		p.setUpMethodFromService(s, m)
		err = p.parseMethodTags(m, ts)
		if err != nil {
			return nil, err
		}
		if err := validateMethod(m); err != nil {
			return nil, err
		}
		s.Methods = append(s.Methods, *m)
	}
	return s, nil
}

func (p *Parser) parseMethod(method *astTypes.Function) (*types.Method, []string, error) {
	m := defaultMethod()
	m.Name = method.Name
	err := p.parseArguments(m, method.Args)
	if err != nil {
		return nil, nil, err
	}
	err = p.parseResults(m, method.Results)
	if err != nil {
		return nil, nil, err
	}
	var ts []string
	ts, m.Docs = cleanComments(method.Docs)
	return m, ts, nil
}

func (p *Parser) parseArguments(m *types.Method, args []astTypes.Variable) error {
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
			a, err := p.parseArgument(arg)
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

func (p *Parser) parseArgument(v astTypes.Variable) (*types.Argument, error) {
	a := defaultArgument()
	var err error
	a.Name = strings.ToUpperFirst(strings.ToCamelCase(v.Name))
	a.Type, err = parseType(v.Type)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (p *Parser) parseResults(m *types.Method, args []astTypes.Variable) error {
	for i, arg := range args {
		if i == len(args)-1 {
			t, ok := arg.Type.(astTypes.TName)
			if !ok || t.TypeName != "error" {
				return fmt.Errorf("last return value in \"%s\" method has to be of type \"error\"", m.Name)
			}
		} else {
			r, err := p.parseField(arg)
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

func (p *Parser) parseField(v astTypes.Variable) (*types.Field, error) {
	f := &types.Field{}
	var err error
	f.Name = strings.ToUpperFirst(strings.ToCamelCase(v.Name))
	f.Type, err = parseType(v.Type)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (p *Parser) parseServiceTags(service *types.Service, tags []string) error {
	for _, tag := range tags {
		builtInTag := extractTag(p.getBuiltInServiceTags(), tag)
		otherTag := extractTag(p.getOtherServiceTags(), tag)
		if len(builtInTag) > len(otherTag) {
			parser := p.builtInServiceTags[strs.ToLower(builtInTag)]
			if err := parser(service, cleanTag(strs.TrimPrefix(tag, "@"+builtInTag))); err != nil {
				return err
			}
			continue
		} else if len(builtInTag) < len(otherTag) {
			parser := p.otherServiceTags[strs.ToLower(otherTag)]
			options := types.TagOptions{}
			if err := parser(*service, options, cleanTag(strs.TrimPrefix(tag, "@"+otherTag))); err != nil {
				return err
			}
			if len(options) > 0 {
				service.OtherOptions[strings.ToLower(otherTag)] = options
			}
			continue
		}
		return fmt.Errorf("invalid service tag \"%s\"", tag)
	}
	return nil
}

func (p *Parser) parseMethodTags(method *types.Method, tags []string) error {
	for _, tag := range tags {
		builtInTag := extractTag(p.getBuiltInMethodTags(), tag)
		otherTag := extractTag(p.getOtherMethodTags(), tag)
		if len(builtInTag) > len(otherTag) {
			parser := p.builtInMethodTags[strs.ToLower(builtInTag)]
			if err := parser(method, cleanTag(strs.TrimPrefix(tag, "@"+builtInTag))); err != nil {
				return err
			}
			continue
		} else if len(builtInTag) < len(otherTag) {
			parser := p.otherMethodTags[strs.ToLower(otherTag)]
			options := types.TagOptions{}
			if err := parser(*method, options, cleanTag(strs.TrimPrefix(tag, "@"+otherTag))); err != nil {
				return err
			}
			if len(options) > 0 {
				method.OtherOptions[strings.ToLower(otherTag)] = options
			}
			continue
		}
		return fmt.Errorf("invalid method tag \"%s\"", tag)
	}
	return nil
}

func (p *Parser) parseParamTags(param *types.Argument, tags []string) error {
	for _, tag := range tags {
		builtInTag := extractTag(p.getBuiltInParamTags(), tag)
		otherTag := extractTag(p.getOtherParamTags(), tag)
		if len(builtInTag) > len(otherTag) {
			parser := p.builtInParamTags[strs.ToLower(builtInTag)]
			if err := parser(param, cleanTag(strs.TrimPrefix(tag, "@"+builtInTag))); err != nil {
				return err
			}
			continue
		} else if len(builtInTag) < len(otherTag) {
			parser := p.otherParamTags[strs.ToLower(otherTag)]
			options := types.TagOptions{}
			if err := parser(*param, options, cleanTag(strs.TrimPrefix(tag, "@"+otherTag))); err != nil {
				return err
			}
			if len(options) > 0 {
				param.OtherOptions[strings.ToLower(otherTag)] = options
			}
			continue
		}
		return fmt.Errorf("invalid param tag \"%s\"", tag)
	}
	return nil
}

func (p *Parser) parseEnums(file *ast.File) (enums []types.Enum, err error) {
	for _, t := range file.Decls {
		if g, ok := t.(*ast.GenDecl); ok && g.Tok.String() == "type" {
			for _, s := range g.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok && fmt.Sprint(ts.Type) == "int" && ts.Name.IsExported() {
					e := types.Enum{}
					e.Name = ts.Name.String()
					e.Cases, err = p.parseEnumCases(file, e.Name)
					if err != nil {
						return nil, err
					}
					enums = append(enums, e)
				}
			}
		}
	}
	return
}

func (p *Parser) parseEnumCases(file *ast.File, enum string) (cases []types.EnumCase, err error) {
	for _, t := range file.Decls {
		if g, ok := t.(*ast.GenDecl); ok && g.Tok.String() == "const" {
			for _, s := range g.Specs {
				if ts, ok := s.(*ast.ValueSpec); ok && fmt.Sprint(ts.Type) == enum {
					if len(ts.Names) != len(ts.Values) {
						return nil, fmt.Errorf("cannot parse '%s' enum values", enum)
					}
					for i := range ts.Names {
						if bv, ok := ts.Values[i].(*ast.BasicLit); ts.Names[i].IsExported() && ok {
							eCase := types.EnumCase{}
							eCase.Name = ts.Names[i].String()
							v, err := strconv.Atoi(bv.Value)
							if err != nil {
								return nil, fmt.Errorf("cannot parse '%s' value '%s' for '%s' enum", enum, eCase.Name, bv.Value)
							}
							eCase.Value = v
							cases = append(cases, eCase)
						} else {
							return nil, fmt.Errorf("cannot parse '%s' enum values", enum)
						}
					}
				}
			}
		}
	}
	return
}

func (p *Parser) parseEntities(ast *astTypes.File) (entities []types.Entity, err error) {
	for _, t := range ast.Structures {
		if !strs.HasSuffix(t.Name, "ArgumentsGroup") {
			e := types.Entity{}
			e.Name = t.Name
			for _, f := range t.Fields {
				field := types.Field{}
				field.Name = f.Name
				field.Tags = f.Tags
				field.Type, err = parseType(f.Variable.Type)
				if err != nil {
					return nil, err
				}
				e.Fields = append(e.Fields, &field)
			}
			entities = append(entities, e)
		}
	}
	return
}

func (p *Parser) setServiceEnums(service *types.Service, enums []types.Enum) {
	for _, e := range enums {
		found := false
		for _, method := range service.Methods {
			for _, arg := range method.Arguments {
				if arg.Type.Name == e.Name {
					arg.Type.Enum = &e
					arg.Type.IsEnum = true
					found = true
				}
			}
			for _, res := range method.Results {
				if res.Type.Name == e.Name {
					res.Type.Enum = &e
					res.Type.IsEnum = true
					found = true
				}
			}
		}
		if found {
			service.Enums = append(service.Enums, e)
		}
	}
}

func (p *Parser) setServiceEntities(service *types.Service, entities []types.Entity) {
	for _, e := range entities {
		found := false
		for _, method := range service.Methods {
			for _, arg := range method.Arguments {
				if arg.Type.Name == e.Name {
					arg.Type.Entity = &e
					arg.Type.IsEntity = true
					found = true
				}
			}
			for _, res := range method.Results {
				if res.Type.Name == e.Name {
					res.Type.Entity = &e
					res.Type.IsEntity = true
					found = true
				}
			}
		}
		if found {
			service.Entities = append(service.Entities, e)
		}
	}
}

func (p *Parser) methodParamsTag(method *types.Method, tag string) error {
	params := strings.SplitS(tag, ",")
	if len(params) != 2 {
		return errInvalidTagFormat("method", "params", tag)
	}
	argNames := strings.SplitS(strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(params[0], "["), "]")), ",")
	args := filterArguments(method.Arguments, func(arg *types.Argument) bool {
		return contains(argNames, strings.ToLowerFirst(arg.Name))
	})
	if len(argNames) != len(args) {
		return fmt.Errorf("invalid arguments '%s' given to params method tag in '%s' method", strs.Join(argNames, ", "), method.Name)
	}
	tags := strings.SplitS(strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(params[1], "("), ")")), ",")
	for _, arg := range args {
		err := p.parseParamTags(arg, tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) serviceGenerateTag(service *types.Service, tag string) error {
	options := strings.SplitS(tag, ",")
	err := p.serviceGenerateFlagsHandler.add(&service.Generate, options...)
	if err != nil {
		if errInvalid, ok := err.(*errGenerateInvalidValue); ok {
			return errInvalidGenerateValue("service", service.Name, "generate", errInvalid.invalidValue)
		}
		return err
	}
	return nil
}

func (p *Parser) serviceGenerateAllTag(service *types.Service, tag string) error {
	options := strings.SplitS(tag, ",")
	err := p.serviceGenerateFlagsHandler.allBut(&service.Generate, options...)
	if err != nil {
		if errInvalid, ok := err.(*errGenerateInvalidValue); ok {
			return errInvalidGenerateValue("service", service.Name, "generate-all", errInvalid.invalidValue)
		}
		return err
	}
	return nil
}

func (p *Parser) methodEnableTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	err := p.methodGenerateFlagsHandler.add(&method.Generate, options...)
	if err != nil {
		if errInvalid, ok := err.(*errGenerateInvalidValue); ok {
			return errInvalidGenerateValue("method", method.Name, "enable", errInvalid.invalidValue)
		}
		return err
	}
	return nil
}

func (p *Parser) methodDisableTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	err := p.methodGenerateFlagsHandler.remove(&method.Generate, options...)
	if err != nil {
		if errInvalid, ok := err.(*errGenerateInvalidValue); ok {
			return errInvalidGenerateValue("method", method.Name, "disable", errInvalid.invalidValue)
		}
		return err
	}
	return nil
}

func (p *Parser) methodEnableAllTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	err := p.methodGenerateFlagsHandler.allBut(&method.Generate, options...)
	if err != nil {
		if errInvalid, ok := err.(*errGenerateInvalidValue); ok {
			return errInvalidGenerateValue("method", method.Name, "enable-all", errInvalid.invalidValue)
		}
		return err
	}
	return nil
}

func (p *Parser) methodDisableAllTag(method *types.Method, tag string) error {
	options := strings.SplitS(tag, ",")
	err := p.methodGenerateFlagsHandler.only(&method.Generate, options...)
	if err != nil {
		if errInvalid, ok := err.(*errGenerateInvalidValue); ok {
			return errInvalidGenerateValue("method", method.Name, "disable-all", errInvalid.invalidValue)
		}
		return err
	}
	return nil
}
