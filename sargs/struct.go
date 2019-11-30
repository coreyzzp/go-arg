package sargs

// argUnit 描述参数或者选项共有的属性
type argUnit struct {
	name   string
	help   string
	defval string
	raw    string
	// 当前单元是否有多个，对于option则对应多个option，对于args，则必须是最后一个
	repeated bool
}

// argOption 描述命令中的一个option
type argOption struct {
	env   string   // 一个option可以通过环境变量来设置
	alias []string // 一个option可以有多个alias，例如-v,--Verbose etc. 但是每个verbose只能对应一个option
	unit  argUnit
}

// argArgs 描述命令行参数
type argArgs struct {
	unit argUnit
}

// subCommands 子命令数组，维护声明的顺序
type subCommands []*argCommand

// subCommands 子命令的map
type subCommandMap map[string]*argCommand

// argCommand 描述一个命令，一个命令可以有多个子命令
type argCommand struct {
	parent    *argCommand // 父命令
	cmds      subCommands // 命令可以有子命令
	cmdMap    subCommandMap
	name      string
	help      string
	example   []string
	options   []*argOption
	optionMap map[string][]*argOption
	args      []*argArgs
	raw       []string
}

// argCli 描述命令行选项
type argCli struct {
	opt *options

	name        string
	author      string
	version     string
	description string
	usage       string
	cmds        subCommands
	cmdMap      subCommandMap
}
