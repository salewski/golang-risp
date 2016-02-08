package runtime

import (
	"fmt"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/util"
)

type RuntimeError struct {
	pos     *lexer.TokenPos
	message string
}

func NewRuntimeError(pos *lexer.TokenPos, format string, data ...interface{}) *RuntimeError {
	return &RuntimeError{
		pos:     pos,
		message: fmt.Sprintf(format, data...),
	}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf(util.Red("runtime error:")+" %s(%d:%d): %s", e.pos.Source.Name(), e.pos.Line, e.pos.Col, e.message)
}
