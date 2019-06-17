package file

import (
	"github.com/wlMalk/goms/parser/types"
)

type (
	serviceGenerator func(file File, service types.Service)
	methodGenerator  func(file File, service types.Service, method types.Method)
)

type serviceGeneratorHandler struct {
	generator  serviceGenerator
	conditions []func(service types.Service) bool
}

type methodGeneratorHandler struct {
	generator  methodGenerator
	conditions []func(service types.Service, method types.Method) bool
	extractor  func(service types.Service) []*types.Method
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
	conditions        []func(service types.Service) bool
	overwrite         bool
	overwriteFunc     func(service types.Service) bool
	merge             bool
	mergeFunc         func(service types.Service) bool
}

func NewSpec(fileType string) Spec {
	f := Spec{}
	f.fileType = "go"
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

func (f Spec) Conditions(conds ...func(service types.Service) bool) Spec {
	f.conditions = conds
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

func (f Spec) Generate(service types.Service) File {
	if !checkServiceConditions(service, f.conditions...) {
		return nil
	}
	file := createFile(f, service, nil)
	applySpecFuncs(file, service, f.beforeFuncs...)
	for _, g := range f.serviceGenerators {
		if !checkServiceConditions(service, g.conditions...) {
			continue
		}
		g.generator(file, service)
	}
	for _, g := range f.methodGenerators {
		methods := service.Methods
		if g.extractor != nil {
			methods = g.extractor(service)
		}
		for _, method := range methods {
			if !checkMethodConditions(service, *method, g.conditions...) {
				continue
			}
			g.generator(file, service, *method)
		}
	}
	return file
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
	//
	// switch f.fileType {
	// case "go":
	// 	file = files.NewGoFile(service.Path, path, name, overwrite, merge)
	// case "proto":
	// 	file = files.NewProtoFile(service.Path, path, name, overwrite, merge)
	// case "text":
	// 	file = files.NewTextFile(service.Path, path, name, f.extension, overwrite, merge)
	// }
	return file
}

func checkServiceConditions(service types.Service, conds ...func(service types.Service) bool) bool {
	for _, cond := range conds {
		if !cond(service) {
			return false
		}
	}
	return true
}

func checkMethodConditions(service types.Service, method types.Method, conds ...func(service types.Service, method types.Method) bool) bool {
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
