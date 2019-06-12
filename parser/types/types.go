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
	HTTP     httpServiceOptions
	GRPC     grpcServiceOptions
	Generate generateServiceOptions
}

type httpServiceOptions struct {
	URIPrefix string
}

type grpcServiceOptions struct {
}

type generateServiceOptions struct {
	Logger           bool
	CircuitBreaking  bool
	RateLimiting     bool
	Recovering       bool
	Caching          bool
	Logging          bool
	Tracing          bool
	ServiceDiscovery bool
	ProtoBuf         bool
	Main             bool
	Validators       bool
	Validating       bool
	Middleware       bool
	MethodStubs      bool
	FrequencyMetric  bool
	LatencyMetric    bool
	CounterMetric    bool
	HTTPServer       bool
	HTTPClient       bool
	GRPCServer       bool
	GRPCClient       bool
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

	Options methodOptions
}

type methodOptions struct {
	HTTP     httpMethodOptions
	GRPC     grpcMethodOptions
	Logging  loggingMethodOptions
	Generate generateMethodOptions
}

type generateMethodOptions struct {
	CircuitBreaking bool
	RateLimiting    bool
	Recovering      bool
	Caching         bool
	Logging         bool
	Validators      bool
	Validating      bool
	Middleware      bool
	MethodStubs     bool
	Tracing         bool
	FrequencyMetric bool
	LatencyMetric   bool
	CounterMetric   bool
	HTTPServer      bool
	HTTPClient      bool
	GRPCServer      bool
	GRPCClient      bool
}

type httpMethodOptions struct {
	Method string
	URI    string
	AbsURI string
}

type grpcMethodOptions struct {
}

type loggingMethodOptions struct {
	IgnoredArguments []string
	IgnoredResults   []string
	LenArguments     []string
	LenResults       []string
	IgnoreError      bool
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

func (t *Type) ProtoBufType() string {
	// TODO
	return ""
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
