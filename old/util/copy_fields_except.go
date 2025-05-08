package util

import (
	"reflect"
)

// CopyFieldsExcept copies exported fields from src -> dst, skipping specified field names.
func CopyFieldsExcept(dst interface{}, src interface{}, skip []string) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()
	dstType := dstVal.Type()

	skipMap := map[string]bool{}
	for _, name := range skip {
		skipMap[name] = true
	}

	for i := 0; i < dstVal.NumField(); i++ {
		field := dstType.Field(i)
		if !skipMap[field.Name] && field.PkgPath == "" { // exported only
			dstVal.Field(i).Set(srcVal.Field(i))
		}
	}
}
