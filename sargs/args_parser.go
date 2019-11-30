package sargs

import "fmt"

// cmdGroup 表示一个命令组，描述的是一个命令，以及其参数与option
type cmdGroup struct {
	raw     []string
	options map[string][]string
	args    []string
}

// func parseOption(token string) (key, value string, ok bool) {
// 	var (
// 		long   bool
// 		keylen = 1
// 	)
// 	if long, ok = isFlag(token); ok {
// 		if long {
// 			keylen = 2
// 		}
// 		sp := strings.Split(token[keylen:], "=")
// 		if len(sp) == 2 {
// 			key = sp[0]
// 			value = sp[1]
// 		} else if len(sp) == 1 {
// 			key = sp[0]
// 		}
// 	}
// 	return
// }

func parseOptionToken(token string) (optKey, optValue string, err error) {
	return
}

func (a *argCli) parseNewOption(cmd *argCommand, optMeta *argOption, optKey, optValue string) (err error) {
	return
}

func (a *argCli) parseNewArg(cmd *argCommand, optMeta *argArgs, arg string) (err error) {
	return
}

// parseArgs 当argcli里面的meta信息都准备好之后才能调用
func (a *argCli) parseArgs(currCmd *argCommand, args []string) (err error) {
	var (
		index        = 0
		isLastOne    = false
		curr, next   string
		currArgIndex = 0
	)

	if err = currCmd.parseFromDefault(); err != nil {
		err = fmt.Errorf("parse from default:%w", err)
		return
	}

	if err = currCmd.parseFromEnv(); err != nil {
		err = fmt.Errorf("parse from env:%w", err)
		return
	}

	for index < len(args) {
		isLastOne = index == len(args)-1
		curr = args[index]
		next = ""
		if !isLastOne {
			next = args[index+1]
		}
		index += 1

		// 当前是option的开头
		if isFlag(curr) {
			var (
				optKey, optValue string
				optMeta          *argOption
			)

			if optKey, optValue, err = parseOptionToken(curr); err != nil {
				err = fmt.Errorf("parse option token %s:%w", curr, err)
				return
			}

			if optMeta = currCmd.findOption(optKey); optMeta == nil {
				err = fmt.Errorf("not found opt %s", curr)
				return
			}

			// 当前option需要value，但是当前option不是使用'='的形式，则认为下一个就是目标
			if optMeta.NeedValue() {
				if optValue == "" {
					if next == "" {
						err = fmt.Errorf("option %s need value", curr)
						return
					}

					if isFlag(next) {
						err = fmt.Errorf("option %s need value, but next symbol %s is flag", curr, next)
						return
					}

					// 直接把下一个符号作为当前opt的value
					optValue = next
					index += 1
				}
			}

			// 尝试解析当前的optKey,optValue
			if err = a.parseNewOption(currCmd, optMeta, optKey, optValue); err != nil {
				err = fmt.Errorf("parse flag %s %s:%w", optKey, optValue, err)
				return
			}
			continue
		}

		// 如果当前有subcmd，则优先判断，如果当前符号是subcmd的开始，那就递归进去处理
		if len(currCmd.subcmds) > 0 {
			if subcmd, ok := currCmd.subcmdMap[curr]; !ok {
				err = newParseError(true, "found unknown subcommand %s,availd subcommand for %s is %s", curr, currCmd, currCmd.subcmds)
				return
			} else if err = a.parseArgs(subcmd, args[index+1:]); err != nil {
				err = fmt.Errorf("parse subcmd %s:%w", subcmd, err)
				return
			}
			continue
		}

		// 当前是一个单独的符号，那就尝试看当前args是否已经解析完毕了
		// meta的合法性校验以及确保了，不会说有subcmd的情况下，还有args
		if currArgIndex < len(currCmd.args) {
			currArgMeta := currCmd.args[currArgIndex]
			if err = a.parseNewArg(currCmd, currArgMeta, curr); err != nil {
				err = fmt.Errorf("parse arg %s:%w", curr, err)
				return
			}

			// 不是repeated才进一步
			if !currArgMeta.IsRepeated() {
				currArgIndex += 1
			}
			continue
		}

		// 这种相当于是非法的args了
		err = newParseError(true, "found invaild args %s", curr)
		return
	}

	return
}
