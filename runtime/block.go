package runtime

import "github.com/raoulvdberge/risp/parser"

type Block struct {
	Nodes []parser.Node
	Scope *Scope
}

func NewBlock(nodes []parser.Node, scope *Scope) *Block {
	return &Block{
		Nodes: nodes,
		Scope: scope,
	}
}
