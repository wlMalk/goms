package parser

import (
	"github.com/wlMalk/goms/parser/types"
)

func defaultService() *types.Service {
	s := &types.Service{}
	s.OtherOptions = types.TagsOptions{}
	return s
}

func defaultMethod() *types.Method {
	m := &types.Method{}
	m.Options.HTTP.Method = "POST"
	m.OtherOptions = types.TagsOptions{}
	return m
}

func defaultArgument() *types.Argument {
	a := &types.Argument{}
	a.Options.HTTP.Origin = "BODY"
	a.OtherOptions = types.TagsOptions{}
	return a
}
