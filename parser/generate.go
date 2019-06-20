package parser

import (
	"fmt"
	"strings"

	"github.com/wlMalk/goms/parser/types"
)

type generateHandler struct {
	allowed types.GenerateList
	groups  map[string][]string
}

func newGenerateHandler() *generateHandler {
	return &generateHandler{
		groups: map[string][]string{},
	}
}

type errGenerateInvalidValue struct {
	invalidValue string
}

func (err *errGenerateInvalidValue) Error() string {
	return fmt.Sprintf("invalid value '%s'", err.invalidValue)
}

func (m *generateHandler) all(g *types.GenerateList) {
	m.copy(g, m.allowed)
}

func (m *generateHandler) allBut(g *types.GenerateList, options ...string) error {
	options = expandGenerateGroups(options, m.groups)
	for _, opt := range options {
		if !m.allowed.Has(opt) {
			return &errGenerateInvalidValue{opt}
		}
	}
	m.all(g)
	g.Remove(options...)
	return nil
}

func (m *generateHandler) add(g *types.GenerateList, options ...string) error {
	options = expandGenerateGroups(options, m.groups)
	for _, opt := range options {
		if !m.allowed.Has(opt) {
			return &errGenerateInvalidValue{opt}
		}
	}
	g.Add(options...)
	return nil
}

func (m *generateHandler) remove(g *types.GenerateList, options ...string) error {
	options = expandGenerateGroups(options, m.groups)
	for _, opt := range options {
		if !m.allowed.Has(opt) {
			return &errGenerateInvalidValue{opt}
		}
	}
	g.Remove(options...)
	return nil
}

func (m *generateHandler) only(g *types.GenerateList, options ...string) error {
	options = expandGenerateGroups(options, m.groups)
	for _, opt := range options {
		if !m.allowed.Has(opt) {
			return &errGenerateInvalidValue{opt}
		}
	}
	g.Empty()
	g.Add(options...)
	return nil
}

func (m *generateHandler) copy(g *types.GenerateList, o types.GenerateList) {
	for _, a := range m.allowed {
		if o.Has(a) {
			g.Add(a)
		}
	}
}

func (m *generateHandler) addAllowed(a ...string) {
	m.allowed.Add(a...)
}

func (m *generateHandler) groupAllowed(name string, a ...string) {
	m.allowed.Add(a...)
	for i := range a {
		a[i] = strings.ToLower(a[i])
	}
	m.groups[strings.ToLower(name)] = a
}
