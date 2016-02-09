package runtime

import "github.com/raoulvdberge/risp/parser"

type MacroHandler func(*Macro, *Block, []parser.Node) (*Value, error)

type Macro struct {
	Types   []string
	Handler MacroHandler
}

func NewMacro(handler MacroHandler, types ...string) *Macro {
	return &Macro{
		Types:   types,
		Handler: handler,
	}
}
