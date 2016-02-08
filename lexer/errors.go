package lexer

import (
	"fmt"
	"github.com/raoulvdberge/risp/util"
)

type SyntaxError struct {
	pos     *TokenPos
	message string
}

func NewSyntaxError(pos *TokenPos, format string, data ...interface{}) *SyntaxError {
	return &SyntaxError{
		pos:     pos,
		message: fmt.Sprintf(format, data...),
	}
}

func (e *SyntaxError) Error() string {
	f := ""

	if e.pos.File != nil {
		f = e.pos.File.Name
	}

	return fmt.Sprintf(util.Red("syntax error:")+" %s(%d:%d): %s", f, e.pos.Line, e.pos.Col, e.message)
}
