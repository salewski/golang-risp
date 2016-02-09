package runtime

import "github.com/raoulvdberge/risp/parser"

func (b *Block) Eval() (*Value, error) {
	var result *Value = Nil

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

	if !b.Scope.HasSymbol(name) {
		return nil, NewRuntimeError(node.Pos(), "unknown symbol '%s'", name)
	}

	return b.Scope.GetSymbol(name), nil
}

func (b *Block) evalList(node *parser.ListNode) (*Value, error) {
	if len(node.Nodes) < 1 {
		return nil, NewRuntimeError(node.Pos(), "invalid list notation: expected a function or macro name")
	}

	if _, ok := node.Nodes[0].(*parser.IdentifierNode); !ok {
		return nil, NewRuntimeError(node.Nodes[0].Pos(), "invalid list notation: expected an identifier")
	}

	nameNode := node.Nodes[0].(*parser.IdentifierNode)

	name := nameNode.Token.Data

	if b.Scope.HasMacro(name) {
		macro := b.Scope.GetMacro(name)
		args := node.Nodes[1:] // omit the macro name

		if len(macro.Types) != len(args) {
			return nil, NewRuntimeError(node.Pos(), "macro '%s' expected %d arguments, got %d", name, len(macro.Types), len(args))
		}

		for i, macroArg := range macro.Types {
			if macroArg != "any" {
				if macroArg != args[i].Name() {
					return nil, NewRuntimeError(node.Pos(), "macro '%s' expected that argument %d should be of type %s, not %s", name, i+1, macroArg, args[i].Name())
				}
			}
		}

		return macro.Handler(macro, args)
	} else {
		if !b.Scope.HasSymbol(name) {
			return nil, NewRuntimeError(node.Pos(), "unknown function '%s'", name)
		}

		value := b.Scope.GetSymbol(name)

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
}
