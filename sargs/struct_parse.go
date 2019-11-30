package sargs

type ParseStyle int

const (
	KParseStyleArgs   ParseStyle = 0
	KParseStyleExtent ParseStyle = 1
)

type options struct {
	style ParseStyle
}

type CliOption interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithParseStyle(s ParseStyle) CliOption {
	return optionFunc(func(o *options) {
		o.style = s
	})
}

func ParseCliArags(cmds []string, opts ...CliOption) (cli *argCli, err error) {
	defopt := &options{}
	for _, o := range opts {
		o.apply(defopt)
	}
	cli = &argCli{}
	return
}
