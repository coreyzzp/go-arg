package sargs

type tagOption struct {
	// ignoreBySpecifyArgsEmpty
	// 通过 arg:"-" 来忽略对应的option如果为false，则认为
	ignoreBySpecifyArgsEmpty bool
}

type options struct {
	tagopt   tagOption
	progName string
	style    ParseStyle
}

func NewDefaultOption() *options {
	return &options{}
}
