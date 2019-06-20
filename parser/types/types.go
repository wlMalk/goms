package types

import (
	"fmt"
)

type Service struct {
	Name       string
	Alias      string
	Docs       []string
	Path       string
	ImportPath string
	Version    Version
	Methods    []Method
	Structs    []*Struct
	Enums      []*Enum

	Options      ServiceOptions
	OtherOptions TagsOptions
	Generate     GenerateList
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

func (v *Version) String() string {
	return v.StringSpecial(".")
}

func (v *Version) StringSpecial(sep string) string {
	if v.Patch != 0 {
		return v.FullStringSpecial(sep)
	} else if v.Minor != 0 {
		return fmt.Sprintf("%d%s%d", v.Major, sep, v.Minor)
	} else {
		return fmt.Sprintf("%d", v.Major)
	}
}

func (v *Version) FullString() string {
	return v.FullStringSpecial(".")
}

func (v *Version) FullStringSpecial(sep string) string {
	return fmt.Sprintf("%d%s%d%s%d", v.Major, sep, v.Minor, sep, v.Patch)
}

type Method struct {
	Name      string
	Alias     string
	Docs      []string
	Arguments []*Argument
	Results   []*Field

	Options      MethodOptions
	OtherOptions TagsOptions
	Generate     GenerateList
}

type Type struct {
	PkgImportPath    string
	Pkg              string
	Name             string
	IsPointer        bool
	IsSlice          bool
	IsVariadic       bool
	IsMap            bool
	IsImport         bool
	IsEntity         bool
	IsStruct         bool
	IsEnum           bool
	IsBuiltin        bool
	IsArgumentsGroup bool
	IsBytes          bool
	Value            *Type
	Struct           *Struct
	Enum             *Enum
	ArgumentsGroup   *ArgumentsGroup
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

func (t *Type) String() (s string) {
	if t.IsVariadic {
		s += "..."
	}
	if t.IsSlice {
		s += "[]"
	}
	if t.IsPointer {
		s += "*"
	}
	if t.IsMap {
		s += "map[" + t.Name + "]" + t.Value.String()
		return
	}
	if t.IsImport {
		s += t.Pkg + "." + t.Name
	} else if !t.IsBytes {
		s += t.Name
	} else {
		s += "[]byte"
	}
	return
}

func (t *Type) GoArgumentType() string {
	return t.String()
}

func (t *Type) GoType() string {
	if t.IsVariadic {
		return "[]" + t.Name
	}
	return t.String()
}

func (t *Type) ProtoBufType() (s string) {
	if t.IsBytes {
		return "bytes"
	}
	if t.IsMap {
		return "map<" + toProtoBufType(t.Name) + ", " + t.Value.ProtoBufType() + ">"
	}
	if t.IsSlice || t.IsVariadic {
		s += "repeated "
	}
	s += toProtoBufType(t.Name)
	return
}

func toProtoBufType(name string) string {
	switch name {
	case "int", "int8", "int16", "int32":
		return "int32"
	case "uint", "uint8", "uint16", "uint32":
		return "uint32"
	case "float64":
		return "double"
	case "float32":
		return "float"
	default:
		return name
	}
}

type ArgumentsGroup struct {
	Name      string
	Docs      []string
	Arguments []*Argument
}

type Struct struct {
	Name   string
	Docs   []string
	Fields []*Field
}

type Enum struct {
	Name  string
	Docs  []string
	Cases []string
}

type Argument struct {
	Name         string
	Docs         []string
	Alias        string
	Type         *Type
	IsOptional   bool
	DefaultValue string

	Options      ArgumentOptions
	OtherOptions TagsOptions
}

type Field struct {
	Name  string
	Docs  []string
	Alias string
	Type  *Type
	Tags  map[string]string
}

type Validator struct {
	Name string
	Args []string
}

type Middleware struct {
	Name string
	Args []string
}
