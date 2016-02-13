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

func SymbolName(namespace string, name string) string {
	if namespace != "" {
		name = namespace + ":" + name
	}

	return name
}

func (b *Block) SymbolName(name string) string {
	return SymbolName(b.Namespace, name)
}
