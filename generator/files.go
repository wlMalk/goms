package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"regexp"
	strs "strings"
	"time"

	"github.com/wlMalk/goms"
	"github.com/wlMalk/goms/generator/strings"
)

type File interface {
	Name() string
	Path() string
	Base() string
	Extension() string
	Overwrite() bool
	Merge() bool
	IsEmpty() bool
	WriteTo(w io.Writer) (n int64, err error)
}

type file struct {
	lines     []string
	name      string
	path      string
	base      string
	extension string
	overwrite bool
	merge     bool
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Path() string {
	return f.path
}

func (f *file) Base() string {
	return f.base
}

func (f *file) Extension() string {
	return f.extension
}

func (f *file) Overwrite() bool {
	return f.overwrite
}

func (f *file) Merge() bool {
	return f.merge
}

func (f *file) IsEmpty() bool {
	return len(f.lines) == 0
}

func (f *file) writeLines(w io.Writer) (n int64, err error) {
	var c int
	for _, l := range f.lines {
		c, err = fmt.Fprintln(w, l)
		n += int64(c)
		if err != nil {
			return
		}
	}
	return
}

type goImportDef struct {
	alias string
	path  string
}

type GoFile struct {
	file
	pkg     string
	imports [][]*goImportDef
}

func (f *GoFile) WriteTo(w io.Writer) (int64, error) {
	f.lines = append(append(f.writePackage(), f.writeImports()...), f.lines...)
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

func generateFileHeader(overwrite bool) (lines []string) {
	lines = append(lines, fmt.Sprintf("// Generated by GoMS v%s (github.com/wlMalk/goms) on %s.", goms.VERSION, time.Now().UTC().Format(time.RFC822)))
	if overwrite {
		lines = append(lines, "// DONT EDIT THIS FILE.")
		lines = append(lines, "// ALL CHANGES WILL BE OVERWRITTEN.")
	} else {
		lines = append(lines, "// You can edit this file.")
	}
	lines = append(lines, "")
	return
}

func (f *GoFile) writePackage() (lines []string) {
	lines = append(lines, generateFileHeader(f.Overwrite())...)
	lines = append(lines, "package "+f.pkg)
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

func (f *GoFile) HasImport(path string) bool {
	for i := range f.imports {
		for _, l := range f.imports[i] {
			if l.path == path {
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

const dnsName string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`

var rxDNSName = regexp.MustCompile(dnsName)

func (f *GoFile) AddImport(alias string, path string) {
	if !f.HasImport(path) {
		if strs.HasPrefix(path, filepath.ToSlash(f.base)) {
			f.imports[1] = append(f.imports[1], &goImportDef{alias: alias, path: path})
		} else if pathParts := filepath.SplitList(path); len(pathParts) > 0 && rxDNSName.MatchString(pathParts[0]) {
			f.imports[2] = append(f.imports[2], &goImportDef{alias: alias, path: path})
		} else {
			f.imports[0] = append(f.imports[0], &goImportDef{alias: strs.TrimSpace(alias), path: path})
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

func (f *GoFile) P(s string) {
	f.lines = append(f.lines, s)
}

func (f *GoFile) C(s string) {
	f.Pf("// %s", s)
}

func (f *GoFile) Cs(s ...string) {
	for _, c := range s {
		f.Pf("// %s", c)
	}
}

func (f *GoFile) Cf(format string, args ...interface{}) {
	f.Pf("// "+format, args...)
}

func (f *GoFile) Pf(format string, args ...interface{}) {
	f.P(fmt.Sprintf(format, args...))
}

func NewGoFile(base string, path string, name string, overwrite bool, merge bool) *GoFile {
	f := &GoFile{}
	f.base = base
	f.path = path
	f.name = name
	f.extension = "go"
	f.pkg = strings.ToSnakeCase(filepath.Base(path))
	f.imports = make([][]*goImportDef, 3)
	f.overwrite = overwrite
	return f
}

type Files []File

func (fs Files) Save() error {
	for _, f := range fs {
		if f.IsEmpty() {
			continue
		}
		fileDir := filepath.Join(f.Base(), f.Path())
		if _, err := os.Stat(fileDir); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(fileDir, 0700)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		filePath := filepath.Join(fileDir, f.Name()+"."+f.Extension())
		_, err := os.Stat(filePath)
		if (err != nil && os.IsNotExist(err)) || (err == nil && f.Overwrite()) {
			file, err := os.Create(filePath)
			if err != nil && os.IsExist(err) {
				err = os.Remove(filePath)
				if err != nil {
					return err
				}
				file, err = os.Create(filePath)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			defer func() {
				err := file.Close()
				if err != nil {
					panic(err)
				}
			}()
			_, err = f.WriteTo(file)
			if err != nil {
				return err
			}
		} else if err == nil && f.Merge() {
			// TODO
		} else if err != nil {
			return err
		}
		// } else if !f.Overwrite() {
		// 	continue
		// }
	}
	return nil
}
