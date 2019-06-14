package files

import (
	"fmt"
	"io"
)

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
