package files

import (
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
