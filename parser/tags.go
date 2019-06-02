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

var serviceTagsParsers = map[string]func(service *types.Service, tag string) error{
	"gen":             parseServiceGenTag,
	"http-URI-prefix": parseServiceHttpUriPrefixTag,
}

var methodTagsParsers = map[string]func(method *types.Method, tag string) error{
	"http-method":  parseMethodHttpMethodTag,
	"http-URI":     parseMethodHttpUriTag,
	"http-abs-URI": parseMethodHttpAbsUriTag,
	"params":       parseMethodParamsTag,
	"logs-ignore":  parseMethodLogsIgnoreTag,
	"alias":        parseMethodAliasTag,
	"validate":     parseMethodValidateTag,
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
			if strs.HasPrefix(tag, serviceTag) {
				found = true
				if err := serviceTagsParsers[serviceTag](service, tag); err != nil {
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
			if strs.HasPrefix(tag, methodTag) {
				found = true
				if err := methodTagsParsers[methodTag](method, tag); err != nil {
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

func parseServiceGenTag(service *types.Service, tag string) error {
	return nil
}

func parseServiceHttpUriPrefixTag(service *types.Service, tag string) error {
	return nil
}

func parseMethodHttpMethodTag(method *types.Method, tag string) error {
	return nil
}

func parseMethodHttpUriTag(method *types.Method, tag string) error {
	return nil
}

func parseMethodHttpAbsUriTag(method *types.Method, tag string) error {
	return nil
}

func parseMethodParamsTag(method *types.Method, tag string) error {
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
