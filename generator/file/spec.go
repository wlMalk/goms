package file

import (
	"strings"

	"github.com/wlMalk/goms/parser/types"
)

type (
	ServiceGenerator        func(file File, service types.Service) error
	MethodGenerator         func(file File, service types.Service, method types.Method) error
	ArgumentsGroupGenerator func(file File, service types.Service, argsGroup types.ArgumentsGroup) error
	EntityGenerator         func(file File, service types.Service, entity types.Entity) error
	EnumGenerator           func(file File, service types.Service, enum types.Enum) error
)

type (
	ServiceCondition        func(service types.Service) bool
	MethodCondition         func(service types.Service, method types.Method) bool
	ArgumentsGroupCondition func(service types.Service, argsGroup types.ArgumentsGroup) bool
	EntityCondition         func(service types.Service, entity types.Entity) bool
	EnumCondition           func(service types.Service, enum types.Enum) bool
)

type (
	MethodsExtractor func(service types.Service) []types.Method
)

type serviceGeneratorHandler struct {
	generator  ServiceGenerator
	conditions []ServiceCondition
}

type methodGeneratorHandler struct {
	generator  MethodGenerator
	conditions []MethodCondition
	extractor  MethodsExtractor
}

type argumentsGroupGeneratorHandler struct {
	generator  ArgumentsGroupGenerator
	conditions []ArgumentsGroupCondition
}

type entityGeneratorHandler struct {
	generator  EntityGenerator
	conditions []EntityCondition
}

type enumGeneratorHandler struct {
	generator  EnumGenerator
	conditions []EnumCondition
}

type SpecBeforeFunc func(file File, service types.Service)
type SpecAfterFunc func(b []byte, service types.Service) ([]byte, error)

type Spec struct {
	path                     string
	pathFunc                 func(service types.Service) string
	name                     string
	nameFunc                 func(service types.Service) string
	fileType                 string
	beforeFuncs              []SpecBeforeFunc
	afterFuncs               []SpecAfterFunc
	serviceGenerators        map[string]serviceGeneratorHandler
	methodGenerators         map[string]methodGeneratorHandler
	argumentsGroupGenerators map[string]argumentsGroupGeneratorHandler
	entityGenerators         map[string]entityGeneratorHandler
	enumGenerators           map[string]enumGeneratorHandler
	conditions               []ServiceCondition
	overwrite                bool
	overwriteFunc            func(service types.Service) bool
	merge                    bool
	mergeFunc                func(service types.Service) bool
}

func NewSpec(fileType string) Spec {
	f := Spec{
		serviceGenerators:        map[string]serviceGeneratorHandler{},
		methodGenerators:         map[string]methodGeneratorHandler{},
		argumentsGroupGenerators: map[string]argumentsGroupGeneratorHandler{},
		entityGenerators:         map[string]entityGeneratorHandler{},
		enumGenerators:           map[string]enumGeneratorHandler{},
	}
	f.fileType = fileType
	return f
}

func (f Spec) Path(path string, pathFunc func(service types.Service) string) Spec {
	f.path = path
	f.pathFunc = pathFunc
	return f
}

func (f Spec) Name(name string, nameFunc func(service types.Service) string) Spec {
	f.name = name
	f.nameFunc = nameFunc
	return f
}

func (f Spec) Overwrite(overwrite bool, overwriteFunc func(service types.Service) bool) Spec {
	f.overwrite = overwrite
	f.overwriteFunc = overwriteFunc
	return f
}

func (f Spec) Merge(merge bool, mergeFunc func(service types.Service) bool) Spec {
	f.merge = merge
	f.mergeFunc = mergeFunc
	return f
}

func (f Spec) Conditions(conds ...ServiceCondition) Spec {
	f.conditions = append(f.conditions, conds...)
	return f
}

func (f Spec) Before(before ...SpecBeforeFunc) Spec {
	f.beforeFuncs = before
	return f
}

func (f Spec) After(after ...SpecAfterFunc) Spec {
	f.afterFuncs = after
	return f
}

func (f Spec) AddServiceGenerator(name string, generator ServiceGenerator, conds ...ServiceCondition) Spec {
	name = strings.ToLower(name)
	g := serviceGeneratorHandler{
		generator:  generator,
		conditions: conds,
	}
	f.serviceGenerators[name] = g
	return f
}

func (f Spec) ServiceGenerator(name string, generator ServiceGenerator) Spec {
	name = strings.ToLower(name)
	g := f.getServiceGenerator(name)
	g.generator = generator
	f.serviceGenerators[name] = g
	return f
}

func (f Spec) ServiceGeneratorConditions(name string, conds ...ServiceCondition) Spec {
	name = strings.ToLower(name)
	g := f.getServiceGenerator(name)
	g.conditions = append(g.conditions, conds...)
	f.serviceGenerators[name] = g
	return f
}

func (f Spec) AddArgumentsGroupGenerator(name string, generator ArgumentsGroupGenerator, conds ...ArgumentsGroupCondition) Spec {
	name = strings.ToLower(name)
	g := argumentsGroupGeneratorHandler{
		generator:  generator,
		conditions: conds,
	}
	f.argumentsGroupGenerators[name] = g
	return f
}

func (f Spec) ArgumentsGroupGenerator(name string, generator ArgumentsGroupGenerator) Spec {
	name = strings.ToLower(name)
	g := f.getArgumentsGroupGenerator(name)
	g.generator = generator
	f.argumentsGroupGenerators[name] = g
	return f
}

func (f Spec) ArgumentsGroupGeneratorConditions(name string, conds ...ArgumentsGroupCondition) Spec {
	name = strings.ToLower(name)
	g := f.getArgumentsGroupGenerator(name)
	g.conditions = append(g.conditions, conds...)
	f.argumentsGroupGenerators[name] = g
	return f
}

func (f Spec) AddEntityGenerator(name string, generator EntityGenerator, conds ...EntityCondition) Spec {
	name = strings.ToLower(name)
	g := entityGeneratorHandler{
		generator:  generator,
		conditions: conds,
	}
	f.entityGenerators[name] = g
	return f
}

func (f Spec) EntityGenerator(name string, generator EntityGenerator) Spec {
	name = strings.ToLower(name)
	g := f.getEntityGenerator(name)
	g.generator = generator
	f.entityGenerators[name] = g
	return f
}

func (f Spec) EntityGeneratorConditions(name string, conds ...EntityCondition) Spec {
	name = strings.ToLower(name)
	g := f.getEntityGenerator(name)
	g.conditions = append(g.conditions, conds...)
	f.entityGenerators[name] = g
	return f
}

func (f Spec) AddEnumGenerator(name string, generator EnumGenerator, conds ...EnumCondition) Spec {
	name = strings.ToLower(name)
	g := enumGeneratorHandler{
		generator:  generator,
		conditions: conds,
	}
	f.enumGenerators[name] = g
	return f
}

func (f Spec) EnumGenerator(name string, generator EnumGenerator) Spec {
	name = strings.ToLower(name)
	g := f.getEnumGenerator(name)
	g.generator = generator
	f.enumGenerators[name] = g
	return f
}

func (f Spec) EnumGeneratorConditions(name string, conds ...EnumCondition) Spec {
	name = strings.ToLower(name)
	g := f.getEnumGenerator(name)
	g.conditions = append(g.conditions, conds...)
	f.enumGenerators[name] = g
	return f
}

func (f Spec) AddMethodGenerator(name string, generator MethodGenerator, extractor MethodsExtractor, conds ...MethodCondition) Spec {
	name = strings.ToLower(name)
	g := methodGeneratorHandler{
		generator:  generator,
		extractor:  extractor,
		conditions: conds,
	}
	f.methodGenerators[name] = g
	return f
}

func (f Spec) MethodGenerator(name string, generator MethodGenerator) Spec {
	name = strings.ToLower(name)
	g := f.getMethodGenerator(name)
	g.generator = generator
	f.methodGenerators[name] = g
	return f
}

func (f Spec) MethodGeneratorConditions(name string, conds ...MethodCondition) Spec {
	name = strings.ToLower(name)
	g := f.getMethodGenerator(name)
	g.conditions = append(g.conditions, conds...)
	f.methodGenerators[name] = g
	return f
}

func (f Spec) MethodGeneratorExtractor(name string, extractor MethodsExtractor) Spec {
	name = strings.ToLower(name)
	g := f.getMethodGenerator(name)
	g.extractor = extractor
	f.methodGenerators[name] = g
	return f
}

func (f Spec) RemoveGenerator(name string) Spec {
	delete(f.serviceGenerators, strings.ToLower(name))
	delete(f.methodGenerators, strings.ToLower(name))
	return f
}

func (f Spec) getServiceGenerator(name string) serviceGeneratorHandler {
	if h, ok := f.serviceGenerators[name]; ok {
		return h
	}
	return serviceGeneratorHandler{}
}

func (f Spec) getArgumentsGroupGenerator(name string) argumentsGroupGeneratorHandler {
	if h, ok := f.argumentsGroupGenerators[name]; ok {
		return h
	}
	return argumentsGroupGeneratorHandler{}
}

func (f Spec) getEntityGenerator(name string) entityGeneratorHandler {
	if h, ok := f.entityGenerators[name]; ok {
		return h
	}
	return entityGeneratorHandler{}
}

func (f Spec) getEnumGenerator(name string) enumGeneratorHandler {
	if h, ok := f.enumGenerators[name]; ok {
		return h
	}
	return enumGeneratorHandler{}
}

func (f Spec) getMethodGenerator(name string) methodGeneratorHandler {
	if h, ok := f.methodGenerators[name]; ok {
		return h
	}
	return methodGeneratorHandler{}
}

func (f Spec) Type() string {
	return f.fileType
}

func (f Spec) Generate(service types.Service, creator Creator) (File, error) {
	if !checkServiceConditions(service, f.conditions...) {
		return nil, nil
	}
	if len(f.serviceGenerators) == 0 &&
		len(f.methodGenerators) == 0 &&
		len(f.entityGenerators) == 0 &&
		len(f.argumentsGroupGenerators) == 0 &&
		len(f.enumGenerators) == 0 {
		return nil, nil
	}
	var err error
	file := createFile(f, service, creator)
	applySpecFuncs(file, service, f.beforeFuncs...)

	err = f.generateService(service, file)
	if err != nil {
		return nil, err
	}
	err = f.generateMethods(service, file)
	if err != nil {
		return nil, err
	}
	err = f.generateEntities(service, file)
	if err != nil {
		return nil, err
	}
	err = f.generateArgumentsGroups(service, file)
	if err != nil {
		return nil, err
	}
	err = f.generateEnums(service, file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f Spec) generateService(service types.Service, file File) (err error) {
	for _, g := range f.serviceGenerators {
		if !checkServiceConditions(service, g.conditions...) {
			continue
		}
		err = g.generator(file, service)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f Spec) generateMethods(service types.Service, file File) (err error) {
	for _, g := range f.methodGenerators {
		methods := service.Methods
		if g.extractor != nil {
			methods = g.extractor(service)
		}
		for _, method := range methods {
			if !checkMethodConditions(service, method, g.conditions...) {
				continue
			}
			err = g.generator(file, service, method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f Spec) generateEntities(service types.Service, file File) (err error) {
	for _, g := range f.entityGenerators {
		entities := service.Entities
		for _, entity := range entities {
			if !checkEntityConditions(service, entity, g.conditions...) {
				continue
			}
			err = g.generator(file, service, entity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f Spec) generateEnums(service types.Service, file File) (err error) {
	for _, g := range f.enumGenerators {
		enums := service.Enums
		for _, enum := range enums {
			if !checkEnumConditions(service, enum, g.conditions...) {
				continue
			}
			err = g.generator(file, service, enum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f Spec) generateArgumentsGroups(service types.Service, file File) (err error) {
	for _, g := range f.argumentsGroupGenerators {
		argsGroups := service.ArgumentsGroups
		for _, argsGroup := range argsGroups {
			if !checkArgumentsGroupConditions(service, argsGroup, g.conditions...) {
				continue
			}
			err = g.generator(file, service, argsGroup)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createFile(f Spec, service types.Service, creator Creator) File {
	var file File
	path := f.path
	if f.pathFunc != nil {
		path = f.pathFunc(service)
	}
	name := f.name
	if f.nameFunc != nil {
		name = f.nameFunc(service)
	}
	overwrite := f.overwrite
	if f.overwriteFunc != nil {
		overwrite = f.overwriteFunc(service)
	}
	merge := f.merge
	if f.mergeFunc != nil {
		merge = f.mergeFunc(service)
	}
	file = creator(service.Path, path, name, overwrite, merge)
	return file
}

func checkServiceConditions(service types.Service, conds ...ServiceCondition) bool {
	for _, cond := range conds {
		if !cond(service) {
			return false
		}
	}
	return true
}

func checkMethodConditions(service types.Service, method types.Method, conds ...MethodCondition) bool {
	for _, cond := range conds {
		if !cond(service, method) {
			return false
		}
	}
	return true
}

func checkEntityConditions(service types.Service, entity types.Entity, conds ...EntityCondition) bool {
	for _, cond := range conds {
		if !cond(service, entity) {
			return false
		}
	}
	return true
}

func checkArgumentsGroupConditions(service types.Service, argsGroup types.ArgumentsGroup, conds ...ArgumentsGroupCondition) bool {
	for _, cond := range conds {
		if !cond(service, argsGroup) {
			return false
		}
	}
	return true
}

func checkEnumConditions(service types.Service, enum types.Enum, conds ...EnumCondition) bool {
	for _, cond := range conds {
		if !cond(service, enum) {
			return false
		}
	}
	return true
}

func applySpecFuncs(file File, service types.Service, funcs ...SpecBeforeFunc) {
	for _, f := range funcs {
		f(file, service)
	}
}
