package sargs

import "fmt"

type ParseError struct {
	cmd *argCommand
	msg string
}

func (p *ParseError) Error() string {
	return p.msg
}

func newParseError(c *argCommand, format string, args ...interface{}) *ParseError {
	return &ParseError{
		cmd: c,
		msg: fmt.Sprintf(format, args...),
	}
}
