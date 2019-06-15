package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/wlMalk/goms/version"
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

func generateFileHeader(overwrite bool) (lines []string) {
	if overwrite {
		lines = append(lines, "Code generated by GoMS; DO NOT EDIT.")
		lines = append(lines, "ALL CHANGES WILL BE OVERWRITTEN.")
	} else {
		lines = append(lines, "Code generated by GoMS; you can edit this file.")
	}
	lines = append(lines, fmt.Sprintf("Generated By: GoMS v%s (github.com/wlMalk/goms)", version.VERSION))
	lines = append(lines, fmt.Sprintf("Generated At: %s", time.Now().UTC().Format(time.RFC822)))
	lines = append(lines, "")
	return
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
		fileName := ""
		if f.Extension() != "" {
			fileName = f.Name() + "." + f.Extension()
		} else {
			fileName = f.Name()
		}
		filePath := filepath.Join(fileDir, fileName)
		_, err := os.Stat(filePath)
		if (err != nil && os.IsNotExist(err)) || (err == nil && f.Overwrite()) {
			buf := new(bytes.Buffer)
			_, err = f.WriteTo(buf)
			if err != nil {
				return err
			}
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
			_, err = io.Copy(file, buf)
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
