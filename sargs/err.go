package sargs

import "fmt"

type ParseError struct {
	needHelp bool
	msg      string
}

func (p *ParseError) Error() string {
	return p.msg
}

func newParseError(needHelp bool, format string, args ...interface{}) *ParseError {
	return &ParseError{
		needHelp: needHelp,
		msg:      fmt.Sprintf(format, args...),
	}
}
