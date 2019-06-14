package files

import (
	"io"
	"path/filepath"

	"github.com/wlMalk/goms/generator/strings"
)

type protoImportDef struct {
	alias string
	path  string
}

type ProtoFile struct {
	TextFile
	Pkg     string
	imports [][]*protoImportDef
}

func NewProtoFile(base string, path string, name string, overwrite bool, merge bool) *ProtoFile {
	f := &ProtoFile{}
	f.base = base
	f.path = path
	f.name = name
	f.extension = "proto"
	f.Pkg = strings.ToSnakeCase(filepath.Base(path))
	f.overwrite = overwrite
	f.merge = merge
	f.CommentFormat("// %s")
	return f
}

func (f *ProtoFile) WriteTo(w io.Writer) (int64, error) {
	lines := f.lines
	f.lines = nil
	f.Cs(generateFileHeader(f.Overwrite())...)
	f.Pf("syntax = \"proto3\";")
	f.Pf("package %s;", f.Pkg)
	f.P("")
	f.lines = append(f.lines, lines...)
	lines = nil
	return f.writeLines(w)

}
