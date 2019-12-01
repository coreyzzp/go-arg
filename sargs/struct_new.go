package sargs

type ParseStyle int

const (
	// 与https://github.com/alexflint/go-arg解析方式兼容，支持其所有功能
	KParseStyleArgs ParseStyle = 0
	// 扩展了一些功能，例如cmdDo等
	KParseStyleExtent ParseStyle = 1
)

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

func NewCli(opts ...CliOption) (cli *argCli) {
	defopt := &options{}
	for _, o := range opts {
		o.apply(defopt)
	}
	// cli.new
	return
}
