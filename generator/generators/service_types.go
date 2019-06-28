package generators

import (
	strs "strings"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func ServiceEntityType(file file.File, service types.Service, entity types.Entity) error {
	entityName := strings.ToUpperFirst(entity.Name)
	file.Pf("type %s struct {", entityName)
	for _, field := range entity.Fields {
		fieldName := strings.ToUpperFirst(field.Name)
		file.Pf("%s %s", fieldName, field.Type.GoType())
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func ServiceArgumentsGroupType(file file.File, service types.Service, argGroup types.ArgumentsGroup) error {
	argGroupName := strings.ToUpperFirst(argGroup.Name)
	file.Pf("type %s struct {", argGroupName)
	for _, arg := range argGroup.Arguments {
		argName := strings.ToUpperFirst(arg.Name)
		file.Pf("%s %s", argName, arg.Type.GoType())
	}
	file.Pf("}")
	file.Pf("")
	return nil
}

func ServiceEnumType(file file.File, service types.Service, enum types.Enum) error {
	enumName := strings.ToUpperFirst(enum.Name)
	file.Pf("type %s int", enumName)
	file.Pf("const(")
	for _, c := range enum.Cases {
		caseName := strs.ToUpper(strings.ToSnakeCase(c.Name))
		file.Pf("%s %s = %d", caseName, enumName, c.Value)
	}
	file.Pf(")")
	file.Pf("")
	return nil
}
