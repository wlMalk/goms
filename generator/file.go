package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/wlMalk/goms/generator/file"
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/version"
)

func GoFileCreator(base string, path string, name string, overwrite bool, merge bool) file.File {
	return files.NewGoFile(base, path, name, overwrite, merge)
}

func ProtoFileCreator(base string, path string, name string, overwrite bool, merge bool) file.File {
	return files.NewProtoFile(base, path, name, overwrite, merge)
}

func TextFileCreator(ext string) file.Creator {
	return file.Creator(func(base string, path string, name string, overwrite bool, merge bool) file.File {
		return files.NewTextFile(base, path, name, ext, overwrite, merge)
	})
}

func generateFileHeader(overwrite bool) (lines []string) {
	if overwrite {
		lines = append(lines, fmt.Sprintf("Code generated by GoMS (v%s); DO NOT EDIT.", version.VERSION))
		lines = append(lines, "ALL CHANGES WILL BE OVERWRITTEN.")
	} else {
		lines = append(lines, fmt.Sprintf("Code generated by GoMS (v%s); You can edit this file.", version.VERSION))
	}
	lines = append(lines, "Generated By: github.com/wlMalk/goms")
	lines = append(lines, fmt.Sprintf("Generated At: %s", time.Now().UTC().Format(time.RFC822)))
	lines = append(lines, "")
	return
}

type Files []file.File

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
			// } else if err == nil && f.Merge() {
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
