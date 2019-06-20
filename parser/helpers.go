package parser

import (
	"fmt"
	strs "strings"
	"unicode/utf8"

	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"

	astTypes "github.com/vetcher/go-astra/types"
)

func (p *Parser) setUpMethodFromService(s *types.Service, m *types.Method) {
	p.methodGenerateFlagsHandler.copy(&m.Generate, s.Generate)
}

func validateMethod(m *types.Method) error {
	return nil
}

func validateArgument(a *types.Argument) error {
	return nil
}

func validateField(f *types.Field) error {
	return nil
}

func getServiceInterfaces(interfaces []astTypes.Interface, serviceName string) (ifaces []*astTypes.Interface) {
	for _, i := range interfaces {
		if isServiceInterface(i.Name, serviceName) {
			ifaces = append(ifaces, &i)
		}
	}
	return
}

func isServiceInterface(name string, serviceName string) bool {
	name = strings.ToUpperFirst(name)
	return strs.HasPrefix(name, serviceName+"Service") ||
		strs.HasPrefix(name, "Service") ||
		strs.HasPrefix(name, serviceName) ||
		strs.HasSuffix(name, serviceName+"Service") ||
		strs.HasSuffix(name, "Service") ||
		strs.HasSuffix(name, serviceName) ||
		strs.HasPrefix(name, serviceName+"Svc") ||
		strs.HasPrefix(name, "Svc") ||
		strs.HasPrefix(name, serviceName) ||
		strs.HasSuffix(name, serviceName+"Svc") ||
		strs.HasSuffix(name, "Svc") ||
		strs.HasSuffix(name, serviceName)
}

func extractServiceVersion(name string, serviceName string) string {
	name = strings.ToUpperFirst(name)
	name = strs.Replace(name, serviceName, "", 1)
	name = strs.Replace(name, "Service", "", 1)
	name = strs.Replace(name, "Svc", "", 1)
	name = strs.TrimLeft(name, "_")
	name = strs.TrimRight(name, "_")
	return name
}

func cleanServiceName(name string) string {
	name = strs.TrimSuffix(name, "_service")
	name = strs.TrimSuffix(name, "_svc")
	name = strings.ToCamelCase(name)
	name = strings.ToUpperFirst(name)
	return name
}

func limitLineLength(str string, length int) []string {
	words := strs.Fields(str)
	var lines []string
	var charCount int
	var line []string
	var count int
	var wordCount int
	for i := 0; i < len(words); {
		wordCount = utf8.RuneCountInString(words[i])
		count = charCount + wordCount + len(line)
		if count <= length {
			line = append(line, words[i])
			charCount += wordCount
			i++
		}
		if count > length || i == len(words) {
			lines = append(lines, strs.Join(line, " "))
			line = nil
			charCount = 0
		}
	}
	return lines
}

func cleanComments(comments []string) (tags []string, docs []string) {
	for i := range comments {
		comments[i] = strs.TrimSpace(comments[i])
		comments[i] = strs.Replace(comments[i], "\n", "", -1)
		comments[i] = strs.Replace(comments[i], "\t", " ", -1)
		comments[i] = strs.TrimPrefix(strs.TrimPrefix(strs.TrimSuffix(comments[i], "*/"), "/*"), "//")
		comments[i] = strs.TrimSpace(comments[i])
	}
	comments = strings.SplitS(strs.Join(comments, " "), " ")
	ignoring := false
	for i := range comments {
		comments[i] = strs.TrimSpace(comments[i])
		if comments[i] == "##" {
			ignoring = true
			continue
		}
		if comments[i][0] == '@' && !ignoring {
			tags = append(tags, comments[i])
		} else if !ignoring && !strs.HasPrefix(comments[i], "##@") {
			docs = append(docs, comments[i])
		}
		ignoring = false
	}
	return strings.SplitS(strs.Join(tags, " "), " "), limitLineLength(strs.Join(docs, " "), 80)
}

func cleanTag(tag string) string {
	return strs.TrimSpace(strs.TrimSuffix(strs.TrimPrefix(strs.TrimSpace(tag), "("), ")"))
}

func errInvalidTagFormat(origin string, tagName string, tag string) error {
	return fmt.Errorf("invalid %s tag '%s': @%s(%s)", origin, tagName, tagName, tag)
}

func errInvalidGenerateValue(origin string, originName string, tagName string, value string) error {
	return fmt.Errorf("invalid value '%s' for %s %s tag in '%s' %s", value, tagName, origin, originName, origin)
}

func expandGenerateGroups(a []string, m map[string][]string) (g []string) {
	for i := range a {
		if n, ok := m[strings.ToLower(a[i])]; ok {
			g = append(g, n...)
		} else {
			g = append(g, strings.ToLower(a[i]))
		}
	}
	return
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

func extractTag(tags []string, tag string) string {
	longest := 0
	for _, t := range tags {
		if len(t) > longest && strs.HasPrefix(strs.ToLower(tag), "@"+t) {
			longest = len(t)
		}
	}
	if longest == 0 {
		return ""
	}
	return tag[1 : longest+1]
}
