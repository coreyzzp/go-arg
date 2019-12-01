package sargs

import (
	"testing"
)

func TestParseOldStyleMeta(t *testing.T) {
	var (
		cli        *argCli
		testTarget struct {
			Verbose   bool
			TestTimes int
		}
	)
	MustParse("prog --help")
}
