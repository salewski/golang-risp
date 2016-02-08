package parser

import (
	"github.com/raoulvdberge/risp/lexer"
)

type Node interface {
	Name() string
	Pos() *lexer.TokenPos
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

type NumberNode struct {
	Token *lexer.Token `json:"token"`
}

func (n *NumberNode) Name() string {
	return "number"
}

func (n *NumberNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
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

type KeywordNode struct {
	Token *lexer.Token `json:"token"`
}

func (n *KeywordNode) Name() string {
	return "keyword"
}

func (n *KeywordNode) Pos() *lexer.TokenPos {
	return n.Token.Pos
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
