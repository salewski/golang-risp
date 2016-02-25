package parser

import (
	"github.com/raoulvdberge/risp/lexer"
)

type Node interface {
	Name() string
	Pos() *lexer.TokenPos
	String() string
}

type StringNode struct {
	Token *lexer.Token `json:"token"`
}

func (n *StringNode) Name() string {
	return "string"
}

func (n *StringNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
}

func (n *StringNode) String() string {
	return n.Token.Data
}

type NumberNode struct {
	Token *lexer.Token `json:"token"`
}

func (n *NumberNode) Name() string {
	return "number"
}

func (n *NumberNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
}

func (n *NumberNode) String() string {
	return n.Token.Data
}

type IdentifierNode struct {
	Token *lexer.Token `json:"token"`
}

func (n *IdentifierNode) Name() string {
	return "identifier"
}

func (n *IdentifierNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
}

func (n *IdentifierNode) String() string {
	return n.Token.Data
}

type KeywordNode struct {
	Token *lexer.Token `json:"token"`
}

func (n *KeywordNode) Name() string {
	return "keyword"
}

func (n *KeywordNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
}

func (n *KeywordNode) String() string {
	return n.Token.Data
}

type ListNode struct {
	OpenToken  *lexer.Token `json:"open"`
	CloseToken *lexer.Token `json:"close"`
	Nodes      []Node       `json:"nodes"`
}

func (n *ListNode) Name() string {
	return "list"
}

func (n *ListNode) Pos() *lexer.TokenPos {
	return n.OpenToken.Pos
}

func (n *ListNode) String() string {
	s := "("

	for i, elem := range n.Nodes {
		s += elem.String()

		if i != len(n.Nodes)-1 {
			s += " "
		}
	}

	s += ")"

	return s
}

type QuoteNode struct {
	Token *lexer.Token `json:"token"`
	Node  Node         `json:"node"`
}

func (n *QuoteNode) Name() string {
	return "quote"
}

func (n *QuoteNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
}

func (n *QuoteNode) String() string {
	return n.Node.String()
}
