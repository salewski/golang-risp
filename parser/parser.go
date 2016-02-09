package parser

import "github.com/raoulvdberge/risp/lexer"

type Parser struct {
	pos    int            `json:"-"`
	Tokens []*lexer.Token `json:"tokens"`
	Nodes  []Node         `json:"nodes"`
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		Tokens: tokens,
	}
}

func (p *Parser) addNode(node Node) {
	p.Nodes = append(p.Nodes, node)
}

func (p *Parser) next() {
	p.pos++
}

func (p *Parser) hasNext() bool {
	return p.pos < len(p.Tokens)
}

func (p *Parser) current() *lexer.Token {
	return p.Tokens[p.pos]
}

func (p *Parser) last() *lexer.Token {
	return p.Tokens[len(p.Tokens)-1]
}

func (p *Parser) nextNode() (Node, error) {
	var node Node

	t := p.current()

	switch {
	case t.IsType(lexer.Identifier):
		node = &IdentifierNode{Token: t}

		p.next()
	case t.IsType(lexer.Keyword):
		node = &KeywordNode{Token: t}

		p.next()
	case t.IsType(lexer.Number):
		node = &NumberNode{Token: t}

		p.next()
	case t.IsType(lexer.String):
		node = &StringNode{Token: t}

		p.next()
	case t.IsTypeAndData(lexer.Separator, "("):
		listNode := &ListNode{
			OpenToken: t,
		}

		depth := 1

		p.next()

		var tokens []*lexer.Token

		for p.hasNext() {
			depth += p.current().DepthModifier()

			if depth == 0 {
				listNode.CloseToken = p.current()

				p.next()

				break
			}

			tokens = append(tokens, p.current())

			p.next()
		}

		if depth != 0 {
			return nil, lexer.NewSyntaxError(p.last().Pos, "unclosed list")
		}

		parser := NewParser(tokens)

		err := parser.Parse()

		if err != nil {
			return nil, err
		}

		listNode.Nodes = parser.Nodes

		node = listNode
	default:
		return nil, lexer.NewSyntaxError(t.Pos, "unexpected token '%s'", t.Data)
	}

	return node, nil
}

func (p *Parser) Parse() error {
	for p.hasNext() {
		node, err := p.nextNode()

		if err != nil {
			return err
		}

		p.addNode(node)
	}

	return nil
}
