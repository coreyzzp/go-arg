package sargs

import (
	"fmt"
	"path/filepath"
	"reflect"
)

var (
	// ErrHelp indicates that -h or --help were provided
	ErrHelp = fmt.Errorf("help requested by user")

	// ErrVersion indicates that --version was provided
	ErrVersion = fmt.Errorf("version requested by user")
)

type orgArgCli struct {
	*argCli
}

func (a *orgArgCli) parseProgName(args []string) {
	if len(args) > 0 {
		if a.opt.progName != "" {
			a.name = filepath.Base(args[0])
		} else {
			a.name = a.opt.progName
		}
	} else {
		a.name = "program"
	}
}

func (a *orgArgCli) parseOrgOneDest(args []string, dest ...interface{}) (err error) {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		err = fmt.Errorf("%s is not a pointer (did you forget an ampersand?)", t)
		return
	}

	// 遍历当前struct的所有field，先解析出来
	for index := 0; index < t.NumField(); index++ {
		// eachField := t.Field(index)
		// fieldTag := eachField.Tag

		// 如果是匿名
	}

	return
}

func (a *orgArgCli) parseOrgStyle(args []string, dest ...interface{}) (err error) {
	a.parseProgName(args)
	for _, dst := range dest {
		if err = a.parseOrgOneDest(args, dst); err != nil {
			return
		}
	}
	return
}
