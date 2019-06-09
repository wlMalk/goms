package generate_service

import (
	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/helpers"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateResponseFile(base string, path string, name string, methods []*types.Method) *files.GoFile {
	file := files.NewGoFile(base, path, name, true, false)
	for _, m := range methods {
		var fields []*types.Field
		for _, res := range m.Results {
			f := &types.Field{}
			f.Name = res.Name
			f.Type = res.Type
			f.Alias = res.Alias
			fields = append(fields, f)
		}
		if len(fields) == 0 {
			continue
		}
		helpers.GenerateExportedStruct(file, m.Name+"Response", fields)
	}
	return file
}
