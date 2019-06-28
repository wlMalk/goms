package file

import (
	"io"
)

type Creator func(base string, path string, name string, overwrite bool, merge bool) File

type File interface {
	Name() string
	Path() string
	Base() string
	Extension() string
	Overwrite() bool
	Merge() bool
	IsEmpty() bool
	WriteTo(w io.Writer) (n int64, err error)
	AddImport(alias string, path ...string)
	HasImport(path ...string) bool
	FormatComments(s ...string) []string
	C(s string)
	Cs(s ...string)
	Cf(format string, args ...interface{})
	P(s string)
	Ps(s ...string)
	Pf(format string, args ...interface{})
}
