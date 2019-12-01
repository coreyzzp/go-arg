package sargs

import (
	"fmt"
	"strings"
)

func (a *argCli) Usage() string {
	return ""
}

func (a *argCli) Help() string {
	return ""
}

func ParseAndExecute(cli *argCli, args []string, dest interface{}) (err error) {
	if err = Parse(cli, args, dest); err != nil {
		err = fmt.Errorf("parse:%w", err)
		return
	}

	// 执行默认动作
	if err = cli.execute(); err != nil {
		err = fmt.Errorf("execute:%w", err)
		return
	}

	return
}

func Parse(cli *argCli, args []string, dest interface{}) (err error) {
	if itf, ok := dest.(interface {
		Version() string
	}); ok {
		cli.version = itf.Version()
	}

	if itf, ok := dest.(interface {
		Description() string
	}); ok {
		cli.description = itf.Description()
	}

	// 解析cli的meta
	if err = cli.parseMeta(args, dest); err != nil {
		err = fmt.Errorf("parse meta:%w", err)
		return
	}

	// 允许动态地增加subcmd
	if itf, ok := dest.(interface {
		SubCmds() []*SubCmd
	}); ok {
		subs := itf.SubCmds()
		for _, s := range subs {
			if err = cli.walkForSubCommand(cli.rootCmd, s); err != nil {
				err = fmt.Errorf("dynamic add command:%w", err)
				return
			}
		}
	}

	// 解析参数
	if err = cli.parseArgs(cli.rootCmd, args); err != nil {
		err = fmt.Errorf("parse args:%w", err)
		return
	}

	return
}

func MustParseArgs(args []string, dest interface{}, opts ...CliOption) {
	cli := NewCli(opts...)
	if err := Parse(cli, args, dest); err != nil {
		panic(err)
	}
}

func MustParseCmdline(cmdline string, dest interface{}, opts ...CliOption) {
	cli := NewCli(opts...)
	args := strings.Split(cmdline, " ")
	if err := Parse(cli, args, dest); err != nil {
		panic(err)
	}
}
