package sargs

type tagOption struct {
	// ignoreBySpecifyArgsEmpty
	// 通过 arg:"-" 来忽略对应的option，如果为true，那么不标记arg的字段也会被当做option处理，这个是goarg的默认行为
	// 如果为false，则一定要有arg字段，才认为当前属于option的一部分，但如果arg的字段里面也有"-"，也会被忽略
	ignoreBySpecifyArgsEmpty bool
}

type walkOption struct {
	// 如果是struct，也递归进去
	walkStructField bool
}

type options struct {
	tagOption
	walkOption
	progName string
	style    ParseStyle
}

func NewDefaultOption() *options {
	return &options{}
}
