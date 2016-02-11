package std

import (
	"fmt"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"math/big"
)

var Symbols = runtime.Symtab{
	"t":       runtime.True,
	"f":       runtime.False,
	"nil":     runtime.Nil,
	"print":   runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrint, "print")),
	"println": runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrintln, "println")),
	"list":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdList, "list")),
	"+":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "+")),
	"-":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "-")),
	"*":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "*")),
	"/":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "/")),
	"=":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdEquals, "=")),
	"!=":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdNotEquals, "!=")),
	">":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, ">")),
	">=":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, ">=")),
	"<":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, "<")),
	"<=":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, "<=")),
	"cat":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdCat, "cat")),
	"and":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdAnd, "and")),
	"or":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdOr, "or")),
	"not":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdNot, "not")),
	"call":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdCall, "call")),
	"range":   runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdRange, "range")),
}

var Macros = runtime.Mactab{
	"defun": runtime.NewMacro(stdDefun, "identifier", "list", "list"),
	"def":   runtime.NewMacro(stdDef, "identifier", "any"),
	"fun":   runtime.NewMacro(stdFun, "list", "list"),
	"for":   runtime.NewMacro(stdFor, "any", "identifier", "list"),
	"while": runtime.NewMacro(stdWhile, "any", "list"),
	"if":    runtime.NewMacro(stdIf, "any", "any"),
	"ifel":  runtime.NewMacro(stdIfel, "any", "any", "any"),
}

func stdDefun(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	name := nodes[0].(*parser.IdentifierNode).Token.Data

	if name == "_" {
		return nil, runtime.NewRuntimeError(nodes[0].Pos(), "it is not allowed to use '_' as a symbol name")
	}

	argNodes := nodes[1].(*parser.ListNode)
	var args []string
	callback := nodes[2].(*parser.ListNode)

	for _, argNode := range argNodes.Nodes {
		ident, ok := argNode.(*parser.IdentifierNode)

		if !ok {
			return nil, runtime.NewRuntimeError(argNode.Pos(), "expected an identifier")
		}

		args = append(args, ident.Token.Data)
	}

	functionBlock := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(block.Scope))
	function := runtime.NewDeclaredFunction(functionBlock, name, args)

	block.Scope.SetSymbol(name, runtime.NewFunctionValue(function))

	return runtime.Nil, nil
}

func stdDef(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	name := nodes[0].(*parser.IdentifierNode).Token.Data
	value, err := block.EvalNode(nodes[1])

	if name == "_" {
		return nil, runtime.NewRuntimeError(nodes[0].Pos(), "it is not allowed to use '_' as a symbol name")
	}

	if err != nil {
		return nil, err
	}

	block.Scope.SetSymbol(name, value)

	return runtime.Nil, nil
}

func stdFun(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	argNodes := nodes[0].(*parser.ListNode)
	var args []string
	callback := nodes[1].(*parser.ListNode)

	for _, argNode := range argNodes.Nodes {
		ident, ok := argNode.(*parser.IdentifierNode)

		if !ok {
			return nil, runtime.NewRuntimeError(argNode.Pos(), "expected an identifier")
		}

		args = append(args, ident.Token.Data)
	}

	functionBlock := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(block.Scope))
	function := runtime.NewLambdaFunction(functionBlock, args)

	return runtime.NewFunctionValue(function), nil
}

func stdFor(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	l, err := block.EvalNode(nodes[0])

	if err != nil {
		return nil, err
	}

	if l.Type != runtime.ListValue {
		return nil, runtime.NewRuntimeError(nodes[0].Pos(), "expected a list to iterate over")
	}

	name := nodes[1].(*parser.IdentifierNode).Token.Data

	callbackBlock := runtime.NewBlock([]parser.Node{nodes[2]}, runtime.NewScope(block.Scope))

	for _, item := range l.List {
		callbackBlock.Scope.SetSymbol(name, item)

		_, err := callbackBlock.Eval()

		if err != nil {
			return nil, err
		}
	}

	return runtime.Nil, nil
}

func stdWhile(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
recheck:
	callback, err := block.EvalNode(nodes[0])

	if err != nil {
		return nil, err
	}

	if callback.Type != runtime.BooleanValue {
		return nil, runtime.NewRuntimeError(nodes[0].Pos(), "expected a boolean")
	}

	if callback.Boolean {
		_, err := block.EvalNode(nodes[1])

		if err != nil {
			return nil, err
		}

		goto recheck
	}

	return runtime.Nil, nil
}

// if and elif are macros because the last argument, the callback
// can't be evaluated if the condition is false.
func stdIf(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	conditionNode := nodes[0]
	condition, err := block.EvalNode(conditionNode)

	if err != nil {
		return nil, err
	}

	if condition.Type != runtime.BooleanValue {
		return nil, runtime.NewRuntimeError(conditionNode.Pos(), "expected a boolean")
	}

	if condition.Boolean == true {
		return block.EvalNode(nodes[1])
	} else {
		return runtime.Nil, nil
	}
}

func stdIfel(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	conditionNode := nodes[0]
	condition, err := block.EvalNode(conditionNode)

	if err != nil {
		return nil, err
	}

	if condition.Type != runtime.BooleanValue {
		return nil, runtime.NewRuntimeError(conditionNode.Pos(), "expected a boolean")
	}

	if condition.Boolean == true {
		return block.EvalNode(nodes[1])
	} else {
		return block.EvalNode(nodes[2])
	}
}

func stdPrint(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Print(arg)
	}

	return runtime.Nil, nil
}

func stdPrintln(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Println(arg)
	}

	return runtime.Nil, nil
}

func stdList(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	l := runtime.NewListValue()

	for _, arg := range context.Args {
		l.List = append(l.List, arg)
	}

	return l, nil
}

func stdMath(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	base := big.NewRat(0, 1)

	var callback func(*big.Rat, *big.Rat) *big.Rat

	switch context.Name {
	case "+":
		callback = base.Add
	case "-":
		callback = base.Sub
	case "*":
		callback = base.Mul
	case "/":
		if context.Args[1].Number.Cmp(base) == 0 {
			return nil, runtime.NewRuntimeError(context.Pos, "division by zero")
		}

		callback = base.Quo
	}

	return runtime.NewNumberValueFromRat(callback(context.Args[0].Number, context.Args[1].Number)), nil
}

func stdMathCmp(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	n1 := context.Args[0].Number
	n2 := context.Args[1].Number

	ok := false

	switch context.Name {
	case ">":
		ok = n1.Cmp(n2) == 1
	case ">=":
		ok = n1.Cmp(n2) >= 0
	case "<":
		ok = n1.Cmp(n2) == -1
	case "<=":
		ok = n1.Cmp(n2) <= 0
	}

	return runtime.BooleanValueFor(ok), nil
}

func stdEquals(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue, runtime.AnyValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(context.Args[0].Equals(context.Args[1])), nil
}

func stdNotEquals(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue, runtime.AnyValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(!context.Args[0].Equals(context.Args[1])), nil
}

func stdCat(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	s := ""

	for _, arg := range context.Args {
		s += arg.String()
	}

	return runtime.NewStringValue(s), nil
}

func stdAnd(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(context.Args[0].Boolean && context.Args[1].Boolean), nil
}

func stdOr(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(context.Args[0].Boolean || context.Args[1].Boolean), nil
}

func stdNot(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(!context.Args[0].Boolean), nil
}

func stdCall(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if len(context.Args) < 1 {
		return nil, runtime.NewRuntimeError(context.Pos, "expected a function")
	}

	function := context.Args[0].Function

	return function.Call(context.Block, context.Args[1:], context.Pos)
}

func stdRange(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	low := context.Args[0].NumberToInt64()
	high := context.Args[1].NumberToInt64()

	if low > high {
		return nil, runtime.NewRuntimeError(context.Pos, "invalid argument, low can't be higher than high (%d > %d)", low, high)
	}

	l := runtime.NewListValue()

	for i := low; i < high; i++ {
		l.List = append(l.List, runtime.NewNumberValueFromInt64(i))
	}

	return l, nil
}
