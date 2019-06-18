package file

import (
	"strings"

	"github.com/wlMalk/goms/parser/types"
)

type (
	ServiceGenerator func(file File, service types.Service) error
	MethodGenerator  func(file File, service types.Service, method types.Method) error
)

type (
	ServiceCondition func(service types.Service) bool
	MethodCondition  func(service types.Service, method types.Method) bool
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

type SpecBeforeFunc func(file File, service types.Service)
type SpecAfterFunc func(b []byte, service types.Service) ([]byte, error)

type Spec struct {
	path              string
	pathFunc          func(service types.Service) string
	name              string
	nameFunc          func(service types.Service) string
	fileType          string
	beforeFuncs       []SpecBeforeFunc
	afterFuncs        []SpecAfterFunc
	serviceGenerators map[string]serviceGeneratorHandler
	methodGenerators  map[string]methodGeneratorHandler
	conditions        []ServiceCondition
	overwrite         bool
	overwriteFunc     func(service types.Service) bool
	merge             bool
	mergeFunc         func(service types.Service) bool
}

func NewSpec(fileType string) Spec {
	f := Spec{
		serviceGenerators: map[string]serviceGeneratorHandler{},
		methodGenerators:  map[string]methodGeneratorHandler{},
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

func (f Spec) getServiceGenerator(name string) serviceGeneratorHandler {
	if h, ok := f.serviceGenerators[name]; ok {
		return h
	}
	return serviceGeneratorHandler{}
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
	var err error
	file := createFile(f, service, creator)
	applySpecFuncs(file, service, f.beforeFuncs...)
	for _, g := range f.serviceGenerators {
		if !checkServiceConditions(service, g.conditions...) {
			continue
		}
		err = g.generator(file, service)
		if err != nil {
			return nil, err
		}
	}
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
				return nil, err
			}
		}
	}
	return file, nil
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

func applySpecFuncs(file File, service types.Service, funcs ...SpecBeforeFunc) {
	for _, f := range funcs {
		f(file, service)
	}
}
