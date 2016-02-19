package list

import (
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
)

var Macros = runtime.Mactab{
	"filter": runtime.NewMacro(listFilter, true, "any", "identifier", "list"),
	"map":    runtime.NewMacro(listMap, true, "any", "identifier", "list"),
	"reduce": runtime.NewMacro(listReduce, true, "any", "identifier", "identifier", "list"),
}

func listFilter(context *runtime.MacroCallContext) (*runtime.Value, error) {
	list, err := context.Block.EvalNode(context.Nodes[0])

	if err != nil {
		return nil, err
	}

	if list.Type != runtime.ListValue {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "expected a list")
	}

	ident := context.Nodes[1].(*parser.IdentifierNode).Token.Data

	callback := context.Nodes[2]

	filteredList := runtime.NewListValue()

	for _, item := range list.List {
		b := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(context.Block.Scope))
		b.Scope.SetSymbolLocally(ident, runtime.NewSymbol(item))

		result, err := b.Eval()

		if err != nil {
			return nil, err
		}

		if result.Type != runtime.BooleanValue {
			return nil, runtime.NewRuntimeError(context.Nodes[2].Pos(), "expected a boolean return value, got %s", result.Type.String())
		}

		if result.Boolean {
			filteredList.List = append(filteredList.List, item)
		}
	}

	return filteredList, nil
}

func listMap(context *runtime.MacroCallContext) (*runtime.Value, error) {
	list, err := context.Block.EvalNode(context.Nodes[0])

	if err != nil {
		return nil, err
	}

	if list.Type != runtime.ListValue {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "expected a list")
	}

	ident := context.Nodes[1].(*parser.IdentifierNode).Token.Data

	callback := context.Nodes[2]

	mappedList := runtime.NewListValue()

	for _, item := range list.List {
		b := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(context.Block.Scope))
		b.Scope.SetSymbolLocally(ident, runtime.NewSymbol(item))

		result, err := b.Eval()

		if err != nil {
			return nil, err
		}

		mappedList.List = append(mappedList.List, result)
	}

	return mappedList, nil
}

func listReduce(context *runtime.MacroCallContext) (*runtime.Value, error) {
	list, err := context.Block.EvalNode(context.Nodes[0])

	if err != nil {
		return nil, err
	}

	if list.Type != runtime.ListValue {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "expected a list")
	}

	if len(list.List) == 0 {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "empty list")
	}

	identLeft := context.Nodes[1].(*parser.IdentifierNode).Token.Data
	identRight := context.Nodes[2].(*parser.IdentifierNode).Token.Data

	callback := context.Nodes[3]

	reduced := list.List[0]

	for _, item := range list.List[1:] {
		b := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(context.Block.Scope))
		b.Scope.SetSymbolLocally(identLeft, runtime.NewSymbol(reduced))
		b.Scope.SetSymbolLocally(identRight, runtime.NewSymbol(item))

		result, err := b.Eval()

		if err != nil {
			return nil, err
		}

		reduced = result
	}

	return reduced, nil
}
