package std

import (
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
)

var Macros = runtime.Mactab{
	"defun": runtime.NewMacro(stdDefun, true, "identifier", "list", "list"),
	"def":   runtime.NewMacro(stdDef, true, "identifier", "any"),
	"fun":   runtime.NewMacro(stdFun, true, "list", "list"),
	"for":   runtime.NewMacro(stdFor, true, "any", "list", "list"),
	"while": runtime.NewMacro(stdWhile, true, "any", "list"),
	"if":    runtime.NewMacro(stdIf, true, "any", "any"),
	"ifel":  runtime.NewMacro(stdIfel, true, "any", "any", "any"),
	"case":  runtime.NewMacro(stdCase, false),
}

func stdDefun(context *runtime.MacroCallContext) (*runtime.Value, error) {
	name := context.Nodes[0].(*parser.IdentifierNode).Token.Data

	if name == "_" {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "disallowed symbol name")
	}

	argNodes := context.Nodes[1].(*parser.ListNode)
	var args []string
	callback := context.Nodes[2].(*parser.ListNode)

	if len(callback.Nodes) == 0 {
		return nil, runtime.NewRuntimeError(callback.Pos(), "empty function body")
	}

	for _, argNode := range argNodes.Nodes {
		ident, ok := argNode.(*parser.IdentifierNode)

		if !ok {
			return nil, runtime.NewRuntimeError(argNode.Pos(), "expected an identifier")
		}

		args = append(args, ident.Token.Data)
	}

	function := runtime.NewDeclaredFunction([]parser.Node{callback}, name, args)

	context.Block.Scope.SetSymbol(name, runtime.NewSymbol(runtime.NewFunctionValue(function)))

	return runtime.Nil, nil
}

func stdDef(context *runtime.MacroCallContext) (*runtime.Value, error) {
	name := context.Nodes[0].(*parser.IdentifierNode).Token.Data
	value, err := context.Block.EvalNode(context.Nodes[1])

	if name == "_" {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "disallowed symbol name")
	}

	if err != nil {
		return nil, err
	}

	context.Block.Scope.SetSymbol(name, runtime.NewSymbol(value))

	return runtime.Nil, nil
}

func stdFun(context *runtime.MacroCallContext) (*runtime.Value, error) {
	argNodes := context.Nodes[0].(*parser.ListNode)
	var args []string
	callback := context.Nodes[1].(*parser.ListNode)

	if len(callback.Nodes) == 0 {
		return nil, runtime.NewRuntimeError(callback.Pos(), "empty function body")
	}

	for _, argNode := range argNodes.Nodes {
		ident, ok := argNode.(*parser.IdentifierNode)

		if !ok {
			return nil, runtime.NewRuntimeError(argNode.Pos(), "expected an identifier")
		}

		args = append(args, ident.Token.Data)
	}

	function := runtime.NewLambdaFunction([]parser.Node{callback}, args)

	return runtime.NewFunctionValue(function), nil
}

func stdFor(context *runtime.MacroCallContext) (*runtime.Value, error) {
	l, err := context.Block.EvalNode(context.Nodes[0])

	if err != nil {
		return nil, err
	}

	if l.Type != runtime.ListValue {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "expected a list to iterate over")
	}

	var args []string

	for _, nameNode := range context.Nodes[1].(*parser.ListNode).Nodes {
		ident, isIdent := nameNode.(*parser.IdentifierNode)

		if isIdent {
			args = append(args, ident.Token.Data)
		} else {
			return nil, runtime.NewRuntimeError(nameNode.Pos(), "expected an identifier")
		}
	}

	if len(args) > 2 {
		return nil, runtime.NewRuntimeError(context.Nodes[1].Pos(), "too many arguments provided")
	}

	callbackBlock := runtime.NewBlock([]parser.Node{context.Nodes[2]}, runtime.NewScope(context.Block.Scope))

	for i, item := range l.List {
		if len(args) >= 1 {
			callbackBlock.Scope.SetSymbol(args[0], runtime.NewSymbol(item))
		}

		if len(args) == 2 {
			callbackBlock.Scope.SetSymbol(args[1], runtime.NewSymbol(runtime.NewNumberValueFromInt64(int64(i))))
		}

		_, err := callbackBlock.Eval()

		if err != nil {
			return nil, err
		}
	}

	return runtime.Nil, nil
}

func stdWhile(context *runtime.MacroCallContext) (*runtime.Value, error) {
recheck:
	callback, err := context.Block.EvalNode(context.Nodes[0])

	if err != nil {
		return nil, err
	}

	if callback.Type != runtime.BooleanValue {
		return nil, runtime.NewRuntimeError(context.Nodes[0].Pos(), "expected a boolean")
	}

	if callback.Boolean {
		_, err := context.Block.EvalNode(context.Nodes[1])

		if err != nil {
			return nil, err
		}

		goto recheck
	}

	return runtime.Nil, nil
}

// if and elif are macros because the last argument, the callback
// can't be evaluated if the condition is false.
func stdIf(context *runtime.MacroCallContext) (*runtime.Value, error) {
	conditionNode := context.Nodes[0]
	condition, err := context.Block.EvalNode(conditionNode)

	if err != nil {
		return nil, err
	}

	if condition.Type != runtime.BooleanValue {
		return nil, runtime.NewRuntimeError(conditionNode.Pos(), "expected a boolean")
	}

	if condition.Boolean == true {
		return context.Block.EvalNode(context.Nodes[1])
	} else {
		return runtime.Nil, nil
	}
}

func stdIfel(context *runtime.MacroCallContext) (*runtime.Value, error) {
	conditionNode := context.Nodes[0]
	condition, err := context.Block.EvalNode(conditionNode)

	if err != nil {
		return nil, err
	}

	if condition.Type != runtime.BooleanValue {
		return nil, runtime.NewRuntimeError(conditionNode.Pos(), "expected a boolean")
	}

	if condition.Boolean == true {
		return context.Block.EvalNode(context.Nodes[1])
	} else {
		return context.Block.EvalNode(context.Nodes[2])
	}
}

type caseElement struct {
	cases    []*runtime.Value
	callback parser.Node
}

func stdCase(context *runtime.MacroCallContext) (*runtime.Value, error) {
	if len(context.Nodes) < 1 {
		return nil, runtime.NewRuntimeError(context.Pos, "missing value to compare to")
	}

	matchNode := context.Nodes[0]
	match, err := context.Block.EvalNode(matchNode)

	if err != nil {
		return nil, err
	}

	if ((len(context.Nodes) - 1) % 2) != 0 { // -1 because we can't count for the match node too
		return nil, runtime.NewRuntimeError(context.Pos, "unbalanced case call")
	}

	var elems []caseElement
	var otherwise parser.Node

	// we begin at 1 because we need to omit the match node
	for i := 1; i < len(context.Nodes); i++ {
		elem := caseElement{}

		list, isList := context.Nodes[i].(*parser.ListNode)

		if isList {
			for _, caseNode := range list.Nodes {
				result, err := context.Block.EvalNode(caseNode)

				if err != nil {
					return nil, err
				}

				elem.cases = append(elem.cases, result)
			}

			i++

			elem.callback = context.Nodes[i]

			elems = append(elems, elem)
		} else {
			ident, isIdent := context.Nodes[i].(*parser.IdentifierNode)

			if isIdent && ident.Token.Data == "_" {
				if otherwise != nil {
					return nil, runtime.NewRuntimeError(ident.Pos(), "match can only have one otherwise case")
				}

				i++

				otherwise = context.Nodes[i]
			} else {
				return nil, runtime.NewRuntimeError(context.Nodes[i].Pos(), "expected a list")
			}
		}
	}

	for _, elem := range elems {
		for _, possibility := range elem.cases {
			if possibility.Equals(match) {
				return context.Block.EvalNode(elem.callback)
			}
		}
	}

	if otherwise != nil {
		return context.Block.EvalNode(otherwise)
	}

	return runtime.Nil, nil
}
