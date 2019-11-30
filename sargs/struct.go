package sargs

import "strings"

// argUnit 描述参数或者选项共有的属性
type argUnit struct {
	tag *argTag

	// 当前单元是否有多个，对于option则对应多个option，对于args，则必须是最后一个
	repeated bool

	// 指向目标，如果是一个repeated，则目标应该是一个slice
	target interface{}

	// 解析的时候赋值
	raw string
}

// argOption 描述命令中的一个option
type argOption argUnit

func (a *argOption) NeedValue() bool {
	// TODO: 看当前是否需要value的值
	return false
}

// argArgs 描述命令行参数
type argArgs argUnit

func (a *argArgs) IsRequired() bool {
	return a.tag.isRequired
}

func (a *argArgs) IsRepeated() bool {
	return a.repeated
}

func (a *argArgs) HasValue() bool {
	return false
}

type subCommands []*argCommand

func (s subCommands) String() string {
	n := make([]string, len(s))
	for i, s := range s {
		n[i] = s.name
	}
	return strings.Join(n, ",")
}

// argCommand 描述一个命令，一个命令可以有多个子命令
type argCommand struct {
	name    string
	example []string
	raw     []string

	// 如果是RootCmd，这指向用户传进来的结构体
	// 如果是子命令，则指向在parse构成中，内部构造的临时对象，在进行命令行解析的时候
	// 会用这个值设置到上层cmd对应的field中
	target interface{}

	args []*argArgs

	options   []*argOption
	optionMap map[string]*argOption

	tag       *argTag     // 只有当前cmd属于另外一个subcmd的时候才存在
	parent    *argCommand // 父命令
	subcmds   subCommands // 子命令
	subcmdMap map[string]*argCommand
}

// 当前赋值了可以执行的命令
func (a *argCommand) executeCommandObject() *argCommand {
	// 看下当前哪个subcommand有赋值
	if len(a.subcmds) > 0 {
		for _, sub := range a.subcmds {
			if sub.target != nil {
				// 递归到最底层的命令
				return sub.executeCommandObject()
			}
		}
	}
	// 没有就是自己
	return a
}

func (a *argCommand) String() string {
	return a.name
}

func (a *argCommand) parseFromEnv() (err error) {
	return
}

func (a *argCommand) parseFromDefault() (err error) {
	return
}

func newArgCommand() *argCommand {
	return &argCommand{
		args:      make([]*argArgs, 0),
		options:   make([]*argOption, 0),
		optionMap: make(map[string]*argOption, 0),
		subcmds:   make([]*argCommand, 0),
		subcmdMap: make(map[string]*argCommand, 0),
	}
}

func (a *argCommand) isLastArgsMutiple() bool {
	arglen := len(a.args)
	if arglen == 0 {
		return false
	}

	return a.args[arglen-1].repeated
}

func (a *argCommand) findOption(optToken string) (opt *argOption) {
	return
}

func (a *argCommand) isOptionExist(tag *argTag) bool {
	if tag.shortAlias != "" {
		if _, ok := a.optionMap[tag.shortAlias]; ok {
			return true
		}
	}
	if tag.longAlias != "" {
		if _, ok := a.optionMap[tag.longAlias]; ok {
			return true
		}
	}
	return false
}

func (a *argCommand) isSubCmdExist(name string) bool {
	if _, ok := a.subcmdMap[name]; ok {
		return true
	}
	return false
}

func (a *argCommand) addArag(arg *argArgs) {
	a.args = append(a.args, arg)
}

func (a *argCommand) addOpt(opt *argOption) {
	a.options = append(a.options, opt)
	if opt.tag.shortAlias != "" {
		a.optionMap[opt.tag.shortAlias] = opt
	}
	if opt.tag.longAlias != "" {
		a.optionMap[opt.tag.longAlias] = opt
	}
}

func (a *argCommand) addSubCmd(cmd *argCommand) {
	a.subcmds = append(a.subcmds, cmd)
	a.subcmdMap[cmd.tag.subcmdName] = cmd
}

// argCli 描述命令行选项
type argCli struct {
	opt *options

	name        string
	author      string
	version     string
	description string
	usage       string
	rootCmd     *argCommand
}
