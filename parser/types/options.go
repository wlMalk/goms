package types

import (
	strs "strings"
)

type GenerateList []string

func (g GenerateList) Has(a ...string) bool {
loop:
	for i := range a {
		for j := range g {
			if strs.ToLower(a[i]) == g[j] {
				continue loop
			}
		}
		return false
	}
	return true
}

func (g GenerateList) HasAny(a ...string) bool {
	if len(a) == 0 {
		return true
	}
	for i := range a {
		for j := range g {
			if strs.ToLower(a[i]) == g[j] {
				return true
			}
		}
	}
	return false
}

func (g GenerateList) HasNone(a ...string) bool {
	for i := range a {
		for j := range g {
			if strs.ToLower(a[i]) == g[j] {
				return false
			}
		}
	}
	return true
}

func (g *GenerateList) Add(a ...string) {
	for i := range a {
		if g.HasNone(a[i]) {
			*g = append(*g, strs.ToLower(a[i]))
		}
	}
}

func (g *GenerateList) Remove(a ...string) {
	for i := range a {
		for j := range *g {
			if strs.ToLower(a[i]) == (*g)[j] {
				*g = append((*g)[:j], (*g)[j+1:]...)
				break
			}
		}
	}
}

func (g *GenerateList) Empty() {
	g.Remove(*g...)
}

type TagOptions map[string]interface{}

type TagsOptions map[string]TagOptions

type ServiceOptions struct {
	HTTP HTTPServiceOptions
	GRPC GRPCServiceOptions
}

type HTTPServiceOptions struct {
	URIPrefix string
}

type GRPCServiceOptions struct {
}

type MethodOptions struct {
	HTTP    HTTPMethodOptions
	GRPC    GRPCMethodOptions
	Logging LoggingMethodOptions
}

type HTTPMethodOptions struct {
	Method string
	URI    string
	AbsURI string
}

type GRPCMethodOptions struct {
}

type LoggingMethodOptions struct {
	IgnoredArguments []string
	IgnoredResults   []string
	LenArguments     []string
	LenResults       []string
	IgnoreError      bool
}

type ArgumentOptions struct {
	HTTP HTTPArgumentOptions
}

type HTTPArgumentOptions struct {
	Origin string
}
