package sargs

import "fmt"

// cmdDoStyle 基于cmdDo的interface
func (a *argCli) cmdDoStyle(cmdObj *argCommand) (err error) {
	if cmdObj.parent == nil {
		if itf, ok := cmdObj.target.(interface {
			BeforeCmdDo() error
		}); ok {
			if err = itf.BeforeCmdDo(); err != nil {
				err = fmt.Errorf("before cmd do:%w", err)
				return
			}
		}

		if itf, ok := cmdObj.target.(interface {
			CmdDo() error
		}); ok {
			if err = itf.CmdDo(); err != nil {
				err = fmt.Errorf("cmd do:%w", err)
				return
			}
		}

		if itf, ok := cmdObj.target.(interface {
			AfterCmdDo() error
		}); ok {
			if err = itf.AfterCmdDo(); err != nil {
				err = fmt.Errorf("after cmd do:%w", err)
				return
			}
		}

		return
	}

	var (
		subCmdName = cmdObj.name
		parentObj  = cmdObj.parent.target
	)

	if parentObj == nil {
		err = fmt.Errorf("parent target is nil")
		return
	}

	if itf, ok := cmdObj.target.(interface {
		BeforeSubCmdDo(subCmdName string, sub interface{}) error
	}); ok {
		if err = itf.BeforeSubCmdDo(subCmdName, parentObj); err != nil {
			err = fmt.Errorf("before subcmd do:%w", err)
			return
		}
	}

	if itf, ok := cmdObj.target.(interface {
		CmdDo(parent interface{}) error
	}); ok {
		if err = itf.CmdDo(parentObj); err != nil {
			err = fmt.Errorf("subcmd do:%w", err)
			return
		}
	}

	if itf, ok := cmdObj.target.(interface {
		AfterSubCmdDo(subCmdName string, sub interface{}) error
	}); ok {
		if err = itf.AfterSubCmdDo(subCmdName, parentObj); err != nil {
			err = fmt.Errorf("after subcmd do:%w", err)
			return
		}
	}

	return
}

// 可以通过这个拿到全局的信息
type ParseOpt struct {
	Name      string      // 当前如果执行的是subcmd，则这个是执行的subcmd的名字
	RootOpt   interface{} // 根选项
	ParentOpt interface{} // 父亲的选项
}

// defaultStyle 默认的style
func (a *argCli) defaultStyle(rootCmd, cmdObj *argCommand) (err error) {
	opt := &ParseOpt{
		Name:    cmdObj.name,
		RootOpt: rootCmd.target,
	}
	if cmdObj.parent != nil {
		opt.ParentOpt = cmdObj.parent.target
	}

	if itf, ok := cmdObj.target.(interface {
		BeforeCliDo(o *ParseOpt) error
	}); ok {
		if err = itf.BeforeCliDo(opt); err != nil {
			err = fmt.Errorf("before cli do:%w", err)
			return
		}
	}

	if itf, ok := cmdObj.target.(interface {
		CliDo(o *ParseOpt) error
	}); ok {
		if err = itf.CliDo(opt); err != nil {
			err = fmt.Errorf("cli do:%w", err)
			return
		}
	}

	if itf, ok := cmdObj.target.(interface {
		AfterCliDo(o *ParseOpt) error
	}); ok {
		if err = itf.AfterCliDo(opt); err != nil {
			err = fmt.Errorf("after cli do:%w", err)
			return
		}
	}

	return
}

// execute 执行parse之后的cmd
func (a *argCli) execute() (err error) {
	var (
		cmdObj *argCommand
	)
	if cmdObj = a.rootCmd.executeCommandObject(); cmdObj == nil {
		err = fmt.Errorf("fail to found cmd obj")
		return
	}

	if cmdObj.target == nil {
		err = fmt.Errorf("command target is inl")
		return
	}

	if a.opt.cmdDoStyle {
		return a.cmdDoStyle(cmdObj)
	} else {
		return a.defaultStyle(a.rootCmd, cmdObj)
	}
}
