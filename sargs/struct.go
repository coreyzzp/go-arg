package sargs

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

// argArgs 描述命令行参数
type argArgs argUnit

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

	tag    *argTag       // 只有当前cmd属于另外一个subcmd的时候才存在
	parent *argCommand   // 父命令
	cmds   []*argCommand // 命令可以有子命令
	cmdMap map[string]*argCommand
}

func newArgCommand() *argCommand {
	return &argCommand{
		args:      make([]*argArgs, 0),
		options:   make([]*argOption, 0),
		optionMap: make(map[string]*argOption, 0),
		cmds:      make([]*argCommand, 0),
		cmdMap:    make(map[string]*argCommand, 0),
	}
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

func (a *argCommand) isSubCmdExist(name string) bool {
	if _, ok := a.cmdMap[name]; ok {
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
	a.cmds = append(a.cmds, cmd)
	a.cmdMap[cmd.tag.subcmdName] = cmd
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
