package strings

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToUpperFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func ToSomeCaseWithSep(sep rune, runeConv func(rune) rune) func(string) string {
	return func(s string) string {
		in := []rune(s)
		n := len(in)
		var runes []rune
		for i, r := range in {
			if isExtendedSpace(r) {
				runes = append(runes, sep)
				continue
			}
			if unicode.IsUpper(r) {
				if i > 0 && sep != runes[i-1] && ((i+1 < n && unicode.IsLower(in[i+1])) || unicode.IsLower(in[i-1])) {
					runes = append(runes, sep)
				}
				r = runeConv(r)
			}
			runes = append(runes, r)
		}
		return string(runes)
	}
}

func isExtendedSpace(r rune) bool {
	return unicode.IsSpace(r) || r == '_' || r == '-' || r == '.'
}

var (
	ToSnakeCase    = ToSomeCaseWithSep('_', unicode.ToLower)
	ToURLSnakeCase = ToSomeCaseWithSep('-', unicode.ToLower)
)

func ToKebabCase(s string) string {
	s = ToSnakeCase(s)
	return strings.Replace(strings.Title(strings.Replace(s, "_", " ", -1)), " ", "-", -1)
}

func ToLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func ToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func IsInStringSlice(what string, where []string) bool {
	for _, item := range where {
		if item == what {
			return true
		}
	}
	return false
}

func FetchTags(strs []string, prefix string) (tags []string) {
	for _, comment := range strs {
		if strings.HasPrefix(comment, prefix) {
			tags = append(tags, strings.Split(strings.Replace(comment[len(prefix):], " ", "", -1), ",")...)
		}
	}
	return
}

func HasTag(strs []string, prefix string) bool {
	return ContainTag(strs, prefix)
}

func ToLower(str string) string {
	if len(str) > 0 && unicode.IsLower(rune(str[0])) {
		return str
	}
	for i := range str {
		if unicode.IsLower(rune(str[i])) {
			// Case, when only first char is upper.
			if i == 1 {
				return strings.ToLower(str[:1]) + str[1:]
			}
			return strings.ToLower(str[:i-1]) + str[i-1:]
		}
	}
	return strings.ToLower(str)
}

// Return last upper char in string or first char if no upper characters founded.
func LastUpperOrFirst(str string) string {
	for i := len(str) - 1; i >= 0; i-- {
		if unicode.IsUpper(rune(str[i])) {
			return string(str[i])
		}
	}
	return string(str[0])
}

// Fetch information from slice of comments (docs).
// Returns appendix of first comment which has tag as prefix.
func FetchMetaInfo(tag string, comments []string) string {
	for _, comment := range comments {
		if len(comment) > len(tag) && strings.HasPrefix(comment, tag) {
			return comment[len(tag)+1:]
		}
	}
	return ""
}

func ContainTag(strs []string, prefix string) bool {
	for _, comment := range strs {
		if strings.HasPrefix(comment, prefix) {
			return true
		}
	}
	return false
}

func LastWordFromName(name string) string {
	lastUpper := strings.LastIndexFunc(name, unicode.IsUpper)
	if lastUpper == -1 {
		lastUpper = 0
	}
	return strings.ToLower(name[lastUpper:])
}

func Reverse(s string) string {
	totalLength := len(s)
	buffer := make([]byte, totalLength)
	for i := 0; i < totalLength; {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		utf8.EncodeRune(buffer[totalLength-i:], r)
	}
	return string(buffer)
}

func CountLeading(s string, seps ...rune) (counter int) {
	for _, r := range s {
		found := false
		for _, sep := range seps {
			if r == sep {
				counter++
				found = true
				break
			} else {
				found = false
			}
		}
		if !found {
			return
		}
	}
	return
}

func Split(s string, sep string, f func(i int, before string, after string) (offset int, to int)) (parts []string) {
	if len(sep) == 0 || len(sep) >= len(s) {
		return nil
	}
	p := 0
	for i := 0; i < len(s)-len(sep)+1; i++ {
		if s[i:i+len(sep)] == sep {
			if offset, to := f(i, s[p:i], s[i+len(sep):]); offset != -1 && to != -1 {
				if len(strings.TrimSpace(s[p+offset:i-to])) > 0 {
					parts = append(parts, s[p+offset:i-to])
				}
				p = i + len(sep)
			}
		}
		if i == len(s)-len(sep) {
			if offset, to := f(i, s[p:], ""); offset != -1 && to != -1 && p+offset < len(s)-to {
				parts = append(parts, s[p+offset:len(s)-to])
			}
		}
	}
	return
}

func SplitS(s string, sep string) []string {
	return Split(s, sep, func(i int, before string, after string) (int, int) {
		if strings.Count(before, "(") != strings.Count(before, ")") ||
			strings.Count(before, "[") != strings.Count(before, "]") ||
			strings.Count(before, "{") != strings.Count(before, "}") ||
			strings.Count(before, "<") != strings.Count(before, ">") {
			return -1, -1
		}
		return CountLeading(before, ' ', '\t', '\n'), CountLeading(Reverse(before), ' ', '\t', '\n')
	})
}
