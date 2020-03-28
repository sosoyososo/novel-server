package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func ModelStructToDocMap(model interface{}, prefix string) string {
	t := reflect.TypeOf(model)
	ft, _ := unwrapType(t)
	return docStructType(ft, prefix, true)
}

func docStructType(st reflect.Type, prefix string, addLevel bool) string {
	if st.Kind() != reflect.Struct {
		return ""
	}
	retStr := ""
	if addLevel {
		retStr += prefix
		retStr += "{\n"
	}
	numV := st.NumField()
	for i := 0; i < numV; i++ {
		f := st.Field(i)
		jsons := strings.Split(f.Tag.Get("json"), ",")
		comment := f.Tag.Get("comment")
		if len(comment) > 0 {
			ft, isSlice := unwrapType(f.Type)
			addLevel := false
			sliceSymbol := ""
			if isSlice {
				sliceSymbol = "[]"
			}
			if ft.Kind() == reflect.Struct {
				prefixAddtion := ""
				if len(jsons[0]) > 0 {
					addLevel = true
					retStr += fmt.Sprintf("\t%v%v%v ://%v\n", prefix, sliceSymbol, jsons[0], comment)
					if addLevel {
						prefixAddtion = "\t"
					}
				}
				retStr += docStructType(ft, prefix+prefixAddtion, addLevel)
			} else if len(jsons) > 0 && len(comment) > 0 && jsons[0] != "-" {
				retStr += fmt.Sprintf("\t%v%v%v : %v\n", prefix, sliceSymbol, jsons[0], comment)
			}
		}
	}
	if addLevel {
		retStr += prefix
		retStr += "}\n"
	}
	return retStr
}

func unwrapType(t reflect.Type) (reflect.Type, bool) {
	isSlice := false
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		if t.Kind() == reflect.Slice {
			isSlice = true
			t = t.Elem()
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	return t, isSlice
}
