package sargs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTagParseIgnore 测试通过option来控制option的可见性
func TestTagParseIgnore(t *testing.T) {
	var (
		atag *argTag
		err  error
	)
	var tbl = []struct {
		tag    string
		expect bool
		opt    *options
	}{
		{`arg:"-"`, true, &options{}}, // 通过arg"-"显式指定忽略
		{``, true, &options{}},        // 没有arg参数，则没有option
		{`arg:"-"`, true, &options{tagOption: tagOption{ignoreBySpecifyArgsEmpty: true}}}, // 如果没有，则也认为有option
		{``, true, &options{tagOption: tagOption{ignoreBySpecifyArgsEmpty: true}}},        // 如果没有，主动指定，认为没有
	}
	for _, tc := range tbl {
		atag, err = newArgTag(tc.tag, tc.opt)
		require.Nil(t, err)
		require.Equalf(t, tc.expect, atag.ignore, "%+v", atag)
	}
}

func TestParseOptions(t *testing.T) {
	var (
		atag *argTag
		err  error
	)

	atag, err = newArgTag(`arg:"-v,--Verbose,required" help:"helpText" default:"1"`, &options{})
	require.Nil(t, err)
	require.Equalf(t, "-v", atag.shortAlias, "%+v", atag)
	require.Equalf(t, "--Verbose", atag.longAlias, "%+v", atag)
	require.Equalf(t, true, atag.isRequired, "%+v", atag)
	require.Equalf(t, "1", atag.defval, "%+v", atag)
	require.Equalf(t, "helpText", atag.help, "%+v", atag)
}

func TestParseArgs(t *testing.T) {
	var (
		atag *argTag
		err  error
	)

	atag, err = newArgTag(`arg:"positional,required" help:"helpText" default:"1"`, &options{})
	require.Nil(t, err)
	require.Equalf(t, true, atag.isPosOpt, "%+v", atag)
	require.Equalf(t, true, atag.isRequired, "%+v", atag)
}

func TestParseSubCommand(t *testing.T) {
	var (
		atag *argTag
		err  error
	)

	atag, err = newArgTag(`arg:"subcommand:up" help:"helpText" `, &options{})
	require.Nil(t, err)
	require.Equalf(t, "up", atag.subcmdName, "%+v", atag)
}
