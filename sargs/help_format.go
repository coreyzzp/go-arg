package sargs

type aligner interface {
	Left() string
	Right() string
}

// 让items输出为对齐的形式
func formatToAlign(items []aligner, middleSpace int) (out []string) {
	return
}

type optHelpText struct {
	Left string
	Desc string
}

// 当前command的help
func (a *argCommand) helpText() string {
	// argsSymbolText := make([]string, 0)
	// for _, opt := range a.options {
	// 	argsSymbolText = append(argsSymbolText, opt.tag)
	// }

	return ""
}
