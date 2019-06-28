package files

import (
	"fmt"
	"io"
)

type TextFile struct {
	file
	commentFormat string
}

func NewTextFile(base string, path string, name string, ext string, overwrite bool, merge bool) *TextFile {
	f := &TextFile{}
	f.base = base
	f.path = path
	f.name = name
	f.extension = ext
	f.overwrite = overwrite
	f.merge = merge
	return f
}

func (f *TextFile) CommentFormat(format string) {
	f.commentFormat = format
}

func (f *TextFile) C(s string) {
	f.C(s)
}

func (f *TextFile) Cs(s ...string) {
	for _, c := range s {
		if c == "" {
			f.P("")
			continue
		}
		f.Pf(f.commentFormat, c)
	}
}

func (f *TextFile) Cf(format string, args ...interface{}) {
	f.C(fmt.Sprintf(format, args...))
}

func (f *TextFile) WriteTo(w io.Writer) (int64, error) {
	lines := f.lines
	f.lines = nil
	if f.commentFormat != "" {
		// f.Cs(generateFileHeader(f.Overwrite())...)
	} else {
		// f.Ps(generateFileHeader(f.Overwrite())...)
	}
	f.lines = append(f.lines, lines...)
	lines = nil
	return f.writeLines(w)
}

func (f *TextFile) FormatComments(cs ...string) (fcs []string) {
	for _, c := range cs {
		if c == "" {
			fcs = append(fcs, "")
			continue
		}
		fcs = append(fcs, fmt.Sprintf(f.commentFormat, c))
	}
	return
}

func (f *TextFile) P(s string) {
	f.lines = append(f.lines, s)
}

func (f *TextFile) Pf(format string, args ...interface{}) {
	f.P(fmt.Sprintf(format, args...))
}

func (f *TextFile) Ps(s ...string) {
	f.lines = append(f.lines, s...)
}

func (f *TextFile) HasImport(path ...string) bool {
	// no op
	return false
}

func (f *TextFile) AddImport(alias string, path ...string) {
	// no op
}
