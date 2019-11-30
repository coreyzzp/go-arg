package sargs

import (
	"reflect"
	"strings"
)

// IsExported
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

// IsTypeExport
func IsTypeExport(t reflect.Type) bool {
	return IsExported(t.Name())
}

// isFlag returns true if a token is a flag such as "-v" or "--user" but not "-" or "--"
func isFlag(s string) bool {
	return strings.HasPrefix(s, "-") && strings.TrimLeft(s, "-") != ""
}
