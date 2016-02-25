package runtime

import (
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
)

type MacroHandler func(*MacroCallContext) (*Value, error)

type Macro struct {
	Types        []string
	Handler      MacroHandler
	typeChecking bool
}

func NewMacro(handler MacroHandler, typeChecking bool, types ...string) *Macro {
	return &Macro{
		Types:        types,
		Handler:      handler,
		typeChecking: typeChecking,
	}
}

type MacroCallContext struct {
	Macro *Macro
	Block *Block
	Nodes []parser.Node
	Pos   *lexer.TokenPos
	Name  string
}
