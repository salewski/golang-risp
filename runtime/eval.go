package runtime

import "github.com/raoulvdberge/risp/parser"

func (b *Block) Eval() (*Value, error) {
	var result *Value

	for _, n := range b.Nodes {
		r, err := b.evalNode(n)

		if err != nil {
			return nil, err
		}

		result = r
	}

	return result, nil
}

func (b *Block) evalNode(node parser.Node) (*Value, error) {
	switch node := node.(type) {
	case *parser.StringNode:
		return b.evalString(node), nil
	case *parser.NumberNode:
		return b.evalNumber(node), nil
	case *parser.KeywordNode:
		return b.evalKeyword(node), nil
	case *parser.IdentifierNode:
		return b.evalIdentifier(node)
	case *parser.ListNode:
		return b.evalList(node)
	default:
		return nil, NewRuntimeError(node.Pos(), "unexpected %s", node.Name())
	}
}

func (b *Block) evalString(node *parser.StringNode) *Value {
	return NewStringValue(node.Token.Data)
}

func (b *Block) evalNumber(node *parser.NumberNode) *Value {
	return NewNumberValueFromString(node.Token.Data)
}

func (b *Block) evalKeyword(node *parser.KeywordNode) *Value {
	return NewKeywordValue(node.Token.Data)
}

func (b *Block) evalIdentifier(node *parser.IdentifierNode) (*Value, error) {
	name := node.Token.Data

	if b.Scope.Get(name) == nil {
		return nil, NewRuntimeError(node.Pos(), "unknown symbol '%s'", name)
	}

	return b.Scope.Get(name), nil
}

func (b *Block) evalList(node *parser.ListNode) (*Value, error) {
	if len(node.Nodes) < 1 {
		return nil, NewRuntimeError(node.Pos(), "malformed function call")
	}

	nameNode := node.Nodes[0].(*parser.IdentifierNode)

	name := nameNode.Token.Data

	if !b.Scope.Has(name) {
		return nil, NewRuntimeError(node.Pos(), "unknown function '%s'", name)
	}

	value := b.Scope.Get(name)

	if value.Type != FunctionValue {
		return nil, NewRuntimeError(node.Pos(), "'%s' is not a function", name)
	}

	var args []*Value

	for _, argNode := range node.Nodes[1:] {
		arg, err := b.evalNode(argNode)

		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return value.Function.Call(b, args, node.Pos())
}
