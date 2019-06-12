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

	Options ServiceOptions
}

type ServiceOptions struct {
	Transports TransportOptions
	HTTP       HTTPServiceOptions
	GRPC       GRPCServiceOptions
	Generate   GenerateServiceOptions
}

type transport struct {
	Server bool
	Client bool
}

type TransportOptions struct {
	HTTP transport
	GRPC transport
}

type HTTPServiceOptions struct {
	URIPrefix string
}

type GRPCServiceOptions struct {
}

type GenerateServiceOptions struct {
	Caching          bool
	Logging          bool
	Metrics          bool
	Tracing          bool
	ServiceDiscovery bool
	ProtoBuf         bool
	Main             bool
	Validators       bool
	Middleware       bool
	MethodStubs      bool
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
	Service   *Service
	Name      string
	Docs      []string
	Tags      []string
	Arguments []*Argument
	Results   []*Field

	Options MethodOptions
}

type MethodOptions struct {
	HTTP           HTTPMethodOptions
	GRPC           GRPCMethodOptions
	LoggingOptions LoggingMethodOptions
	Caching        bool
	Logging        bool
	Validator      bool
	Middleware     bool
	MethodStubs    bool
	Tracing        bool
	Metrics        bool
	Transports     TransportOptions
}

type HTTPMethodOptions struct {
	Method string
	URI    string
	AbsURI string
}

type GRPCMethodOptions struct {
}

type LoggingMethodOptions struct {
	Logging          bool
	IgnoredArguments []string
	IgnoredResults   []string
	LenArguments     []string
	LenResults       []string
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

	Options ArgumentOptions
}

type ArgumentOptions struct {
	HTTP HTTPArgumentOptions
}

type HTTPArgumentOptions struct {
	Origin string
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
