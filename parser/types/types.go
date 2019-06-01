package types

import (
	"fmt"
)

type Service struct {
	Name       string
	Docs       []string
	Tags       []string
	Path       string
	ImportPath string
	Version    Version
	Methods    []*Method
	Structs    []*Struct
	Enums      []*Enum
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
	if v.Patch != 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else if v.Minor != 0 {
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
	} else {
		return fmt.Sprintf("%d", v.Major)
	}
}

type Method struct {
	Name      string
	Docs      []string
	Tags      []string
	Arguments []*Argument
	Results   []*Field
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
