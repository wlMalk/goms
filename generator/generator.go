package generator

import (
	"strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/parser/types"
)

type GeneratorOption func(generator *Generator)

type Generator struct {
	creators map[string]file.Creator
	specs    map[string]file.Spec
}

func (g *Generator) AddCreator(fileType string, creator file.Creator) {
	g.creators[strings.ToLower(fileType)] = creator
}

func (g *Generator) AddSpec(name string, spec file.Spec) {
	g.specs[strings.ToLower(name)] = spec
}

func (g *Generator) RemoveSpec(name string) {
	delete(g.specs, strings.ToLower(name))
}

func (g *Generator) RemoveGenerator(name string, generatorName string) {
	g.AddSpec(name, g.GetSpec(name).RemoveGenerator(generatorName))
}

func (g *Generator) AddServiceGenerator(name string, generatorName string, generator file.ServiceGenerator) {
	g.AddSpec(name, g.GetSpec(name).AddServiceGenerator(generatorName, generator))
}

func (g *Generator) AddServiceGeneratorWithConditions(name string, generatorName string, generator file.ServiceGenerator, conds ...file.ServiceCondition) {
	g.AddSpec(name, g.GetSpec(name).AddServiceGenerator(generatorName, generator, conds...))
}

func (g *Generator) AddMethodGenerator(name string, generatorName string, generator file.MethodGenerator) {
	g.AddSpec(name, g.GetSpec(name).AddMethodGenerator(generatorName, generator, nil))
}

func (g *Generator) AddMethodGeneratorWithConditions(name string, generatorName string, generator file.MethodGenerator, conds ...file.MethodCondition) {
	g.AddSpec(name, g.GetSpec(name).AddMethodGenerator(generatorName, generator, nil, conds...))
}

func (g *Generator) AddMethodGeneratorWithExtractor(name string, generatorName string, generator file.MethodGenerator, extractor file.MethodsExtractor) {
	g.AddSpec(name, g.GetSpec(name).AddMethodGenerator(generatorName, generator, extractor))
}

func (g *Generator) AddMethodGeneratorWithExtractorAndConditions(name string, generatorName string, generator file.MethodGenerator, extractor file.MethodsExtractor, conds ...file.MethodCondition) {
	g.AddSpec(name, g.GetSpec(name).AddMethodGenerator(generatorName, generator, extractor, conds...))
}

func (g *Generator) AddEntityGenerator(name string, generatorName string, generator file.EntityGenerator) {
	g.AddSpec(name, g.GetSpec(name).AddEntityGenerator(generatorName, generator))
}

func (g *Generator) AddEntityGeneratorWithConditions(name string, generatorName string, generator file.EntityGenerator, conds ...file.EntityCondition) {
	g.AddSpec(name, g.GetSpec(name).AddEntityGenerator(generatorName, generator, conds...))
}

func (g *Generator) AddArgumentsGroupGenerator(name string, generatorName string, generator file.ArgumentsGroupGenerator) {
	g.AddSpec(name, g.GetSpec(name).AddArgumentsGroupGenerator(generatorName, generator))
}

func (g *Generator) AddArgumentsGroupGeneratorWithConditions(name string, generatorName string, generator file.ArgumentsGroupGenerator, conds ...file.ArgumentsGroupCondition) {
	g.AddSpec(name, g.GetSpec(name).AddArgumentsGroupGenerator(generatorName, generator, conds...))
}

func (g *Generator) AddEnumGenerator(name string, generatorName string, generator file.EnumGenerator) {
	g.AddSpec(name, g.GetSpec(name).AddEnumGenerator(generatorName, generator))
}

func (g *Generator) AddEnumGeneratorWithConditions(name string, generatorName string, generator file.EnumGenerator, conds ...file.EnumCondition) {
	g.AddSpec(name, g.GetSpec(name).AddEnumGenerator(generatorName, generator, conds...))
}

func (g *Generator) GetSpec(name string) file.Spec {
	return g.specs[strings.ToLower(name)]
}

func (g *Generator) Generate(service types.Service) (files Files, err error) {
	for _, s := range g.specs {
		file, err := s.Generate(service, g.creators[strings.ToLower(s.Type())])
		if err != nil {
			return nil, err
		}
		if file != nil {
			files = append(files, file)
		}
	}
	return
}
