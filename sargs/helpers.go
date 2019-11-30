package sargs

import (
	"reflect"
	"strings"
)

func IsExported(name string) bool {
	if len(name) == 0 {
		return false
	}
	typeNameFirstLetter := string(name[0])
	if typeNameFirstLetter == strings.ToLower(typeNameFirstLetter) { // 当前field unexport
		return false
	}
	return true
}

func IsTypeExport(t reflect.Type) bool {
	return IsExported(t.Name())
}

func IsOption(token string) (long, ok bool) {
	switch {
	case strings.HasPrefix(token, "--"):
		if len(token) > 3 && token[2] != '-' {
			ok = true
			long = true
		}
	case strings.HasPrefix(token, "-"):
		if len(token) > 2 {
			ok = true
		}
	}
	return
}
