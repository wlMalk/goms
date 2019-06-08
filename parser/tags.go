package parser

import (
	"fmt"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

var serviceTags = []string{
	"gen",
	"http-URI-prefix",
}

var methodTags = []string{
	"http-method",
	"http-URI",
	"http-abs-URI",
	"params",
	"logs-ignore",
	"alias",
	"validate",
}

var paramTags = []string{
	"http-origin",
}

var serviceTagsParsers = map[string]func(service *types.Service, tag string) error{
	"gen":             parseServiceGenTag,
	"http-URI-prefix": parseServiceHttpUriPrefixTag,
}

var methodTagsParsers = map[string]func(method *types.Method, tag string) error{
	"http-method":  parseMethodHttpMethodTag,
	"http-URI":     parseMethodHttpUriTag,
	"http-abs-URI": parseMethodHttpAbsUriTag,
	"params":       parseMethodParamsTag,
	"results":      parseMethodResultsTag,
	"logs-ignore":  parseMethodLogsIgnoreTag,
	"alias":        parseMethodAliasTag,
	"validate":     parseMethodValidateTag,
}

var paramTagsParsers = map[string]func(argument *types.Argument, tag string) error{
	"http-origin": parseParamHttpOriginTag,
}

func limitLineLength(str string, length int) []string {
	words := strs.Fields(str)
	var lines []string
	charCount := 0
	var line []string
	for i := 0; i < len(words); {
		count := charCount + len(words[i]) + len(line) - 1
		if count <= length {
			line = append(line, words[i])
			charCount += len(words[i])
			i++
		}
		if count > length || i == len(words) {
			lines = append(lines, strs.Join(line, " "))
			line = []string{}
			charCount = 0
		}
	}
	return lines
}

func cleanComments(comments []string) (tags []string, docs []string) {
	for i := range comments {
		comments[i] = strs.TrimSpace(strs.TrimPrefix(comments[i], "//"))
		if comments[i][0] == '@' {
			tags = append(tags, comments[i])
		} else {
			docs = append(docs, comments[i])
		}
	}
	return strings.SplitS(strs.Join(tags, " "), " "), limitLineLength(strs.Join(docs, " "), 80)
}

func parseServiceTags(service *types.Service, tags []string) error {
	for _, tag := range tags {
		found := false
		for _, serviceTag := range serviceTags {
			if strs.HasPrefix(tag, "@"+serviceTag) {
				found = true
				if err := serviceTagsParsers[serviceTag](service, cleanTag(strs.TrimPrefix(tag, "@"+serviceTag))); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid service tag \"%s\"", tag)
		}
	}
	return nil
}

func parseMethodTags(method *types.Method, tags []string) error {
	for _, tag := range tags {
		found := false
		for _, methodTag := range methodTags {
			if strs.HasPrefix(tag, "@"+methodTag) {
				found = true
				if err := methodTagsParsers[methodTag](method, cleanTag(strs.TrimPrefix(tag, "@"+methodTag))); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid method tag \"%s\"", tag)
		}
	}
	return nil
}

func parseParamTags(arg *types.Argument, tags []string) error {
	for _, tag := range tags {
		found := false
		for _, paramTag := range paramTags {
			if strs.HasPrefix(tag, "@"+paramTag) {
				found = true
				if err := paramTagsParsers[paramTag](arg, cleanTag(strs.TrimPrefix(tag, "@"+paramTag))); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid param tag \"%s\"", tag)
		}
	}
	return nil
}

func cleanTag(tag string) string {
	return strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(strs.TrimSpace(tag), "("), ")"))
}

func parseServiceGenTag(service *types.Service, tag string) error {
	return nil
}

func parseServiceHttpUriPrefixTag(service *types.Service, tag string) error {
	return nil
}

func parseMethodHttpMethodTag(method *types.Method, tag string) error {
	httpMethod := strs.ToUpper(tag)
	if httpMethod != "POST" && httpMethod != "GET" && httpMethod != "PUT" && httpMethod != "DELETE" && httpMethod != "OPTIONS" && httpMethod != "HEAD" {
		return fmt.Errorf("invalid http-method value '%s'", tag)
	}
	method.Options.HTTP.Method = httpMethod
	return nil
}

func parseMethodHttpUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.URI = tag
	return nil
}

func parseMethodHttpAbsUriTag(method *types.Method, tag string) error {
	method.Options.HTTP.AbsURI = tag
	return nil
}

func errInvalidTagFormat(origin string, tagName string, tag string) error {
	return fmt.Errorf("invalid %s tag '%s': @%s(%s)", origin, tagName, tagName, tag)
}

func parseMethodParamsTag(method *types.Method, tag string) error {
	params := strings.SplitS(tag, ",")
	if len(params) != 2 {
		return errInvalidTagFormat("method", "params", tag)
	}
	argNames := strings.SplitS(strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(params[0], "["), "]")), ",")
	args := filterArguments(method.Arguments, func(arg *types.Argument) bool {
		return contains(argNames, strings.ToLowerFirst(arg.Name))
	})
	tags := strings.SplitS(strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(params[1], "("), ")")), ",")
	for _, arg := range args {
		err := parseParamTags(arg, tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseMethodResultsTag(method *types.Method, tag string) error {
	return nil
}

func parseMethodLogsIgnoreTag(method *types.Method, tag string) error {
	return nil
}

func parseMethodAliasTag(method *types.Method, tag string) error {
	return nil
}

func parseMethodValidateTag(method *types.Method, tag string) error {
	return nil
}

func parseParamHttpOriginTag(arg *types.Argument, tag string) error {
	origin := strs.ToUpper(tag)
	if origin != "BODY" && origin != "HEADER" && origin != "QUERY" && origin != "PATH" {
		return fmt.Errorf("invalid http-origin value '%s'", tag)
	}
	arg.Options.HTTP.Origin = origin
	return nil
}

func filterArguments(args []*types.Argument, f func(*types.Argument) bool) (filtered []*types.Argument) {
	for _, arg := range args {
		if f(arg) {
			filtered = append(filtered, arg)
		}
	}
	return
}

func filterFields(fields []*types.Field, f func(*types.Field) bool) (filtered []*types.Field) {
	for _, field := range fields {
		if f(field) {
			filtered = append(filtered, field)
		}
	}
	return
}
