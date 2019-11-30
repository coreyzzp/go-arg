package sargs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStructTag(t *testing.T) {
	var (
		atag *argTag
		err  error
	)

	// 通过arg"-"显式指定忽略
	atag, err = newArgTag(`arg:"-"`, &options{})
	require.Nil(t, err)
	require.Equalf(t, true, atag.ignore, "%+v", atag)

	// 如果没有，则也认为有option
	atag, err = newArgTag(``, &options{})
	require.Nil(t, err)
	require.Equalf(t, false, atag.ignore, "%+v", atag)

	// atag, err = newArgTag(`help:"hello"`, &options{tagopt: tagOption{ignoreBySpecifyArgsEmpty: true}})
	// require.Nil(t, err)
	// require.Equal(t, true, atag.ignore)
}
