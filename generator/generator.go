package generator

import "github.com/wlMalk/goms/generator/file"

type GeneratorOption func(generator *Generator)

type Generator struct {
	creators map[string]file.Creator
	specs    map[string]file.Spec
}
