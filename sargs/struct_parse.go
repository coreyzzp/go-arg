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

// 用于输出日志
type filedInfo struct {
	index  int
	parent reflect.Type
	entry  reflect.StructField
}

func (f *filedInfo) String() string {
	return fmt.Sprintf("%s[%d].%s", f.parent.Name(), f.index, f.entry.Name)
}

func (o *orgArgCli) walk(cmd *argCommand, dest interface{}) (err error) {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		err = fmt.Errorf("target should pe a pointer")
		return
	}
	v = v.Elem()  // 指向对应结构体
	t := v.Type() // 遍历对应结构体类型

	// 遍历当前struct的所有field，先解析出来
	for index := 0; index < t.NumField(); index++ {
		var (
			tag        *argTag
			eachField  = t.Field(index)
			fieldTag   = eachField.Tag
			fieldType  = eachField.Type
			fieldKind  = fieldType.Kind()
			fieldValue = v.Field(index)
			f          = &filedInfo{
				index:  index,
				parent: t,
				entry:  t.Field(index),
			}
			isSlice = false
		)

		// 解析tag，得到反射信息
		if tag, err = newArgTag(string(fieldTag), o.opt); err != nil {
			err = fmt.Errorf("walk for %s:%w", f, err)
			return
		}

		if tag.ignore {
			continue
		}

		// 如果是非匿名的，需要看是不是exported的field，如果非exported的，那也跳过
		// 默认情况下，忽略匿名的field
		if !eachField.Anonymous && !IsExported(eachField.Name) {
			err = fmt.Errorf("unexported filed but has tag for %s :%w", f, err)
			return
		}

		// 子命令支持
		if tag.tagType == KTagTypeSubCmd {
			if fieldKind != reflect.Ptr || fieldType.Elem().Kind() != reflect.Struct {
				err = fmt.Errorf("subcmd type should be a ptr to struct for %s :%w", f, err)
				return
			}
			if cmd.isSubCmdExist(tag.subcmdName) {
				err = fmt.Errorf("sbcmd %s exist in %s", tag.subcmdName, f)
				return
			}
			newCmd := newArgCommand()
			newCmd.name = tag.subcmdName
			newCmd.parent = cmd
			newCmd.tag = tag
			newCmd.target = reflect.New(fieldType.Elem())

			if err = o.walk(newCmd, newCmd.target); err != nil {
				err = fmt.Errorf("walk subcmd at %s:%w", f, err)
				return
			}

			cmd.addSubCmd(newCmd)
			return
		}

		// 有一些不支持的就不处理了
		switch fieldKind {
		case reflect.Ptr, reflect.Map,
			reflect.Uintptr, reflect.Complex64, reflect.Complex128,
			reflect.Array, reflect.Chan, reflect.Func,
			reflect.Interface, reflect.UnsafePointer:
			err = fmt.Errorf("not support type %s", f)
			return
		case reflect.Slice:
			isSlice = true
			// 如果是参数，对于slice必须要是最后一个，即要看下之前是否已经加过这种类型的参数了
			if tag.tagType == KTagTypeArgs && cmd.isLastArgsMutiple() {
				err = fmt.Errorf("[]args should be the last one %s", f)
				return
			}
			// 如果是option，则要看下args是不是之前以及添加过了
			if tag.tagType == KTagTypeOption && cmd.isOptionExist(tag) {
				err = fmt.Errorf("tag option alias exist %s", f)
				return
			}
		case reflect.Struct:
			// 对于无论是否匿名都可以递归进去
			if !o.opt.walkStructField {
				err = fmt.Errorf("not handle struct field for %s", f)
				return
			}
			// 如果要处理，则递归进去
			if err = o.walk(cmd, v.Field(index).Addr()); err != nil {
				err = fmt.Errorf("walk embed struct for %s:%w", f, err)
				return
			}
		default:
		}

		// ok，准备构造参数
		newUnit := &argUnit{}
		newUnit.tag = tag
		newUnit.repeated = isSlice
		newUnit.target = fieldValue.Addr()

		switch {
		case tag.tagType == KTagTypeOption:
			cmd.addOpt((*argOption)(newUnit))
		case tag.tagType == KTagTypeArgs:
			cmd.addArag((*argArgs)(newUnit))
		default:
			err = fmt.Errorf("uknown tag type %v", tag.tagType)
			return
		}
	}

	return
}

func (o *orgArgCli) parseArgs(args []string) (err error) {
	return
}

func (a *orgArgCli) parseOrgStyle(args []string, dest interface{}) (err error) {
	a.parseProgName(args)
	rootCmd := newArgCommand()
	rootCmd.target = dest
	if err = a.walk(rootCmd, dest); err != nil {
		err = fmt.Errorf("parse cmd meta fail:%w", err)
		return
	}
	a.rootCmd = rootCmd
	return
}
