package sargs

// argUnit 描述参数或者选项共有的属性
type argUnit struct {
	tag    *argTag
	name   string
	help   string
	defval string
	raw    string
	// 当前单元是否有多个，对于option则对应多个option，对于args，则必须是最后一个
	repeated bool
	target   interface{} // 指向目标，如果是一个repeated，则目标应该是一个slice
}

// argOption 描述命令中的一个option
type argOption struct {
	argUnit
}

// argArgs 描述命令行参数
type argArgs struct {
	argUnit
}

// argCommand 描述一个命令，一个命令可以有多个子命令
type argCommand struct {
	tag     *argTag
	name    string
	help    string
	example []string
	raw     []string
	target  interface{} // 指向目标，如果是一个repeated，则目标应该是一个slice

	args []*argArgs

	options   []*argOption
	optionMap map[string]*argOption

	parent *argCommand   // 父命令
	cmds   []*argCommand // 命令可以有子命令
	cmdMap map[string]*argCommand
}

func (a *argCommand) isLastArgsMutiple() bool {
	arglen := len(a.args)
	if arglen == 0 {
		return false
	}

	return a.args[arglen-1].repeated
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

func (a *argCommand) addArag() {

}

func (a *argCommand) addOpt() {

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
