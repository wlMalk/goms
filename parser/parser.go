package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wlMalk/goms/parser/tags"
	"github.com/wlMalk/goms/parser/types"
)

var isValidTagName = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString

type (
	serviceTagParser func(service *types.Service, tag string) error
	methodTagParser  func(method *types.Method, tag string) error
	paramTagParser   func(arg *types.Argument, tag string) error
)

type (
	ServiceTagParser func(service types.Service, options types.TagOptions, tag string) error
	MethodTagParser  func(method types.Method, options types.TagOptions, tag string) error
	ParamTagParser   func(arg types.Argument, options types.TagOptions, tag string) error
)

type Parser struct {
	builtInServiceTags map[string]serviceTagParser
	otherServiceTags   map[string]ServiceTagParser
	builtInMethodTags  map[string]methodTagParser
	otherMethodTags    map[string]MethodTagParser
	builtInParamTags   map[string]paramTagParser
	otherParamTags     map[string]ParamTagParser

	serviceGenerateFlagsHandler *generateHandler
	methodGenerateFlagsHandler  *generateHandler
}

type ParserOption func(parser *Parser)

func BuiltInServiceTagsParsers(parser *Parser) {
	parser.registerServiceTagParser("name", tags.ServiceNameTag)
	parser.registerServiceTagParser("transports", tags.ServiceTransportsTag)
	parser.registerServiceTagParser("metrics", tags.ServiceMetricsTag)
	parser.registerServiceTagParser("http-URI-prefix", tags.ServiceHTTPUriPrefixTag)
}

func BuiltInMethodTagsParsers(parser *Parser) {
	parser.registerMethodTagParser("name", tags.MethodNameTag)
	parser.registerMethodTagParser("transports", tags.MethodTransportsTag)
	parser.registerMethodTagParser("metrics", tags.MethodMetricsTag)
	parser.registerMethodTagParser("http-method", tags.MethodHTTPMethodTag)
	parser.registerMethodTagParser("http-URI", tags.MethodHTTPUriTag)
	parser.registerMethodTagParser("http-abs-URI", tags.MethodHTTPAbsUriTag)
	parser.registerMethodTagParser("logs-ignore", tags.MethodLogsIgnoreTag)
	parser.registerMethodTagParser("logs-len", tags.MethodLogsLenTag)
	parser.registerMethodTagParser("alias", tags.MethodAliasTag)
}

func BuiltInParamTagsParsers(parser *Parser) {
	parser.registerParamTagParser("http-origin", tags.ParamHTTPOriginTag)
}

func New(opts ...ParserOption) *Parser {
	p := &Parser{
		builtInServiceTags:          map[string]serviceTagParser{},
		otherServiceTags:            map[string]ServiceTagParser{},
		builtInMethodTags:           map[string]methodTagParser{},
		otherMethodTags:             map[string]MethodTagParser{},
		builtInParamTags:            map[string]paramTagParser{},
		otherParamTags:              map[string]ParamTagParser{},
		serviceGenerateFlagsHandler: newGenerateHandler(),
		methodGenerateFlagsHandler:  newGenerateHandler(),
	}
	p.builtInMethodTags["params"] = p.methodParamsTag

	p.builtInServiceTags["generate-all"] = p.serviceGenerateAllTag
	p.builtInServiceTags["generate"] = p.serviceGenerateTag

	p.builtInMethodTags["disable"] = p.methodDisableTag
	p.builtInMethodTags["enable"] = p.methodEnableTag
	p.builtInMethodTags["disable-all"] = p.methodDisableAllTag
	p.builtInMethodTags["enable-all"] = p.methodEnableAllTag

	for _, opt := range opts {
		opt(p)
	}
	return p
}

func Default(opts ...ParserOption) *Parser {
	opts = append([]ParserOption{BuiltInServiceTagsParsers, BuiltInMethodTagsParsers, BuiltInParamTagsParsers}, opts...)
	return New(opts...)
}

func (p *Parser) RegisterServiceTagParser(name string, parser ServiceTagParser) error {
	if !isValidTagName(name) {
		return fmt.Errorf("tag name '%s' is invalid", name)
	}
	if _, ok := p.builtInServiceTags[strings.ToLower(name)]; ok {
		return fmt.Errorf("service tag with name '%s' already exists", name)
	}
	if _, ok := p.otherServiceTags[strings.ToLower(name)]; ok {
		return fmt.Errorf("service tag with name '%s' already exists", name)
	}
	p.otherServiceTags[strings.ToLower(name)] = parser
	return nil
}

func (p *Parser) MustRegisterServiceTagParser(name string, parser ServiceTagParser) {

}

func (p *Parser) registerServiceTagParser(name string, parser serviceTagParser) {
	p.builtInServiceTags[strings.ToLower(name)] = parser
}

func (p *Parser) RegisterMethodTagParser(name string, parser MethodTagParser) error {
	if !isValidTagName(name) {
		return fmt.Errorf("tag name '%s' is invalid", name)
	}
	if _, ok := p.builtInMethodTags[strings.ToLower(name)]; ok {
		return fmt.Errorf("method tag with name '%s' already exists", name)
	}
	if _, ok := p.otherMethodTags[strings.ToLower(name)]; ok {
		return fmt.Errorf("method tag with name '%s' already exists", name)
	}
	p.otherMethodTags[strings.ToLower(name)] = parser
	return nil
}

func (p *Parser) MustRegisterMethodTagParser(name string, parser MethodTagParser) {

}

func (p *Parser) registerMethodTagParser(name string, parser methodTagParser) {
	p.builtInMethodTags[strings.ToLower(name)] = parser
}

func (p *Parser) RegisterParamTagParser(name string, parser ParamTagParser) error {
	if !isValidTagName(name) {
		return fmt.Errorf("tag name '%s' is invalid", name)
	}
	if _, ok := p.builtInParamTags[strings.ToLower(name)]; ok {
		return fmt.Errorf("param tag with name '%s' already exists", name)
	}
	if _, ok := p.otherParamTags[strings.ToLower(name)]; ok {
		return fmt.Errorf("param tag with name '%s' already exists", name)
	}
	p.otherParamTags[strings.ToLower(name)] = parser
	return nil
}

func (p *Parser) MustRegisterParamTagParser(name string, parser ParamTagParser) {
}

func (p *Parser) registerParamTagParser(name string, parser paramTagParser) {
	p.builtInParamTags[strings.ToLower(name)] = parser
}

func (p *Parser) getBuiltInServiceTags() (tags []string) {
	for k := range p.builtInServiceTags {
		tags = append(tags, k)
	}
	return
}

func (p *Parser) getOtherServiceTags() (tags []string) {
	for k := range p.otherServiceTags {
		tags = append(tags, k)
	}
	return
}

func (p *Parser) getBuiltInMethodTags() (tags []string) {
	for k := range p.builtInMethodTags {
		tags = append(tags, k)
	}
	return
}

func (p *Parser) getOtherMethodTags() (tags []string) {
	for k := range p.otherMethodTags {
		tags = append(tags, k)
	}
	return
}

func (p *Parser) getBuiltInParamTags() (tags []string) {
	for k := range p.builtInParamTags {
		tags = append(tags, k)
	}
	return
}

func (p *Parser) getOtherParamTags() (tags []string) {
	for k := range p.otherParamTags {
		tags = append(tags, k)
	}
	return
}
