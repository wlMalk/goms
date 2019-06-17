package files

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"path/filepath"
	"regexp"
	strs "strings"

	"github.com/wlMalk/goms/generator/strings"
)

type goImportDef struct {
	alias string
	path  string
}

type GoFile struct {
	file
	Pkg     string
	imports [][]*goImportDef
}

func NewGoFile(base string, path string, name string, overwrite bool, merge bool) *GoFile {
	f := &GoFile{}
	f.base = base
	f.path = path
	f.name = name
	f.extension = "go"
	f.Pkg = strings.ToSnakeCase(filepath.Base(path))
	f.imports = make([][]*goImportDef, 3)
	f.overwrite = overwrite
	f.merge = merge
	return f
}

func (f *GoFile) WriteTo(w io.Writer) (int64, error) {
	lines := f.lines
	f.lines = nil
	// f.Cs(generateFileHeader(f.Overwrite())...)
	f.lines = append(f.lines, f.writePackage()...)
	f.lines = append(f.lines, f.writeImports()...)
	f.lines = append(f.lines, lines...)
	lines = nil
	buf := new(bytes.Buffer)
	_, err := f.writeLines(buf)
	if err != nil {
		return 0, err
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, err
	}
	buf.Reset()
	return io.Copy(w, bytes.NewReader(b))
}

func (f *GoFile) writePackage() (lines []string) {
	lines = append(lines, "package "+f.Pkg)
	lines = append(lines, "")
	return
}

func (f *GoFile) writeImports() (lines []string) {
	hasImports := false
	for i := range f.imports {
		if len(f.imports[i]) > 0 {
			hasImports = true
			break
		}
	}
	if !hasImports {
		return
	}
	lines = append(lines, "import (")
	for i := range f.imports {
		for _, l := range f.imports[i] {
			if l.alias == "" {
				lines = append(lines, "\""+l.path+"\"")
			} else {
				lines = append(lines, l.alias+" \""+l.path+"\"")
			}
		}
		lines = append(lines, "")
	}
	lines = append(lines, ")")
	return
}

func (f *GoFile) HasImport(path ...string) bool {
	for i := range path {
		path[i] = strs.Trim(path[i], "/")
	}
	iPath := strs.Join(path, "/")
	for i := range f.imports {
		for _, l := range f.imports[i] {
			if l.path == iPath {
				return true
			}
		}
	}
	return false
}

func (f *GoFile) getImport(path string) *goImportDef {
	for i := range f.imports {
		for _, l := range f.imports[i] {
			if l.path == path {
				return l
			}
		}
	}
	return nil
}

const dnsName string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})+[\._]?$`

var rxDNSName = regexp.MustCompile(dnsName)

func (f *GoFile) AddImport(alias string, path ...string) {
	if !f.HasImport(path...) {
		for i := range path {
			path[i] = strs.Trim(path[i], "/")
		}
		if strs.HasSuffix(filepath.ToSlash(f.base), path[0]) {
			f.imports[1] = append(f.imports[1], &goImportDef{alias: alias, path: strs.Join(path, "/")})
		} else if pathParts := strs.Split(strs.Join(path, "/"), "/"); len(pathParts) > 0 && rxDNSName.MatchString(pathParts[0]) {
			f.imports[2] = append(f.imports[2], &goImportDef{alias: alias, path: strs.Join(path, "/")})
		} else {
			f.imports[0] = append(f.imports[0], &goImportDef{alias: strs.TrimSpace(alias), path: strs.Join(path, "/")})
		}
	}
}

func (f *GoFile) I(path string) string {
	f.AddImport("", path)
	i := f.getImport(path)
	if i.alias == "" {
		return strings.ToSnakeCase(filepath.Base(i.path))
	}
	return i.alias
}

func (f *GoFile) C(s string) {
	f.Cs(s)
}

func (f *GoFile) Cs(s ...string) {
	for _, c := range s {
		if c == "" {
			f.P("")
			continue
		}
		f.Pf("// %s", c)
	}
}

func (f *GoFile) Cf(format string, args ...interface{}) {
	f.Pf("// "+format, args...)
}

func (f *GoFile) P(s string) {
	f.lines = append(f.lines, s)
}

func (f *GoFile) Ps(s ...string) {
	f.lines = append(f.lines, s...)
}

func (f *GoFile) Pf(format string, args ...interface{}) {
	f.P(fmt.Sprintf(format, args...))
}
