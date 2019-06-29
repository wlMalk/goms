package parser

import (
	"fmt"

	"github.com/wlMalk/goms/parser/types"

	astTypes "github.com/vetcher/go-astra/types"
)

func acceptableBuiltin(t string) bool {
	return contains([]string{
		"bool",
		"int64",
		"int32",
		"int16",
		"int8",
		"int",
		"uint64",
		"uint32",
		"uint16",
		"uint8",
		"uint",
		"float64",
		"float32",
		"byte",
		"rune",
		"string",
	}, t)
}

func parseTNameType(t *types.Type, typ astTypes.TName) bool {
	if astTypes.IsBuiltin(typ) && acceptableBuiltin(typ.TypeName) {
		t.Name = typ.TypeName
		t.IsBuiltin = true
	} else if !astTypes.IsBuiltin(typ) {
		t.Name = typ.TypeName
	} else {
		return false
	}
	return true
}

func parseTImportType(t *types.Type, typ astTypes.TImport) bool {
	Ttyp, ok := typ.Next.(astTypes.TName)
	if !ok {
		return false
	}
	t.IsImport = true
	t.Pkg = typ.Import.Name
	t.PkgImportPath = typ.Import.Package
	t.Name = Ttyp.TypeName
	return true
}

func parseTPointerType(t *types.Type, typ astTypes.TPointer) bool {
	t.IsPointer = true
	switch Ttyp := typ.Next.(type) {
	case astTypes.TImport:
		return parseTImportType(t, Ttyp)
	case astTypes.TName:
		return parseTNameType(t, Ttyp) && !astTypes.IsBuiltin(Ttyp)
	default:
		return false
	}
}

func parseTMapType(t *types.Type, typ astTypes.TMap) bool {
	t.IsMap = true
	if v, ok := typ.Key.(astTypes.TName); ok && astTypes.IsBuiltin(typ.Key) &&
		acceptableBuiltin(v.TypeName) && !contains([]string{"float64", "float32"}, v.TypeName) {
		t.Name = v.TypeName
		t.Value = &types.Type{}
		if !parseRepeatedValueType(t.Value, typ.Value) {
			return false
		}
		return true
	}
	return false
}

func parseRepeatedValueType(t *types.Type, typ astTypes.Type) bool {
	switch Ttyp := typ.(type) {
	case astTypes.TImport:
		return parseTImportType(t, Ttyp)
	case astTypes.TPointer:
		return parseTPointerType(t, Ttyp)
	case astTypes.TName:
		return parseTNameType(t, Ttyp)
	case astTypes.TArray:
		if !Ttyp.IsSlice || !astTypes.IsBuiltin(Ttyp.Next) || Ttyp.Next.(astTypes.TName).TypeName != "byte" {
			return false
		}
		t.IsBytes = true
		t.Name = "byte"
		return true
	default:
		return false
	}
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

func parseType(typ astTypes.Type) (*types.Type, error) {
	t := &types.Type{}
	var err error = fmt.Errorf("invalid type %s", typ.String())
	switch Ttyp := typ.(type) {
	case astTypes.TEllipsis:
		t.IsVariadic = true
		if !parseRepeatedValueType(t, Ttyp.Next) {
			return nil, err
		}
	case astTypes.TArray:
		if !Ttyp.IsSlice {
			return nil, err
		}
		t.IsSlice = true
		if !parseRepeatedValueType(t, Ttyp.Next) {
			return nil, err
		}
	case astTypes.TMap:
		if !parseTMapType(t, Ttyp) {
			return nil, err
		}
	case astTypes.TPointer:
		if !parseTPointerType(t, Ttyp) {
			return nil, err
		}
	case astTypes.TName:
		if !parseTNameType(t, Ttyp) {
			return nil, err
		}
	case astTypes.TImport:
		if !parseTImportType(t, Ttyp) {
			return nil, err
		}
	default:
		return nil, err
	}

	return t, nil
}
