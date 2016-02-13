package runtime

import "github.com/raoulvdberge/risp/parser"

type Block struct {
	Nodes     []parser.Node
	Scope     *Scope
	Namespace string
}

func NewBlock(nodes []parser.Node, scope *Scope) *Block {
	return &Block{
		Nodes:     nodes,
		Scope:     scope,
		Namespace: "",
	}
}

func (b *Block) SymbolName(name string) string {
	n := name

	if b.Namespace != "" {
		n = b.Namespace + ":" + n
	}

	return n
}
