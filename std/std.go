package std

import (
	"fmt"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/util"
	"math"
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
	"%":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPow, "%")),
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
	"pass":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPass, "pass")),
	"load":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdLoad, "load")),
	"sqrt":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "sqrt")),
	"sin":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "sin")),
	"cos":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "cos")),
	"tan":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "tan")),
	"ceil":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "ceil")),
	"floor":   runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "floor")),
	"abs":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "abs")),
	"log":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "log")),
	"log10":   runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "log10")),
	"pow":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPow, "pow")),
	"deg2rad": runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdDeg2Rad, "deg2rad")),
	"rad2deg": runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdRad2Deg, "rad2deg")),
	"pi":      runtime.NewNumberValueFromFloat64(math.Pi),
	"e":       runtime.NewNumberValueFromFloat64(math.E),
}

var Macros = runtime.Mactab{
	"defun": runtime.NewMacro(stdDefun, "identifier", "list", "list"),
	"def":   runtime.NewMacro(stdDef, "identifier", "any"),
	"fun":   runtime.NewMacro(stdFun, "list", "list"),
	"for":   runtime.NewMacro(stdFor, "any", "list", "list"),
	"while": runtime.NewMacro(stdWhile, "any", "list"),
	"if":    runtime.NewMacro(stdIf, "any", "any"),
	"ifel":  runtime.NewMacro(stdIfel, "any", "any", "any"),
}

func stdDefun(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	name := nodes[0].(*parser.IdentifierNode).Token.Data

	if name == "_" {
		return nil, runtime.NewRuntimeError(nodes[0].Pos(), "disallowed symbol name")
	}

	argNodes := nodes[1].(*parser.ListNode)
	var args []string
	callback := nodes[2].(*parser.ListNode)

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

	functionBlock := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(block.Scope))
	function := runtime.NewDeclaredFunction(functionBlock, name, args)

	block.Scope.SetSymbol(name, runtime.NewFunctionValue(function))

	return runtime.Nil, nil
}

func stdDef(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	name := nodes[0].(*parser.IdentifierNode).Token.Data
	value, err := block.EvalNode(nodes[1])

	if name == "_" {
		return nil, runtime.NewRuntimeError(nodes[0].Pos(), "disallowed symbol name")
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

	var args []string

	for _, nameNode := range nodes[1].(*parser.ListNode).Nodes {
		ident, isIdent := nameNode.(*parser.IdentifierNode)

		if isIdent {
			args = append(args, ident.Token.Data)
		} else {
			return nil, runtime.NewRuntimeError(nameNode.Pos(), "expected an identifier")
		}
	}

	if len(args) > 2 {
		return nil, runtime.NewRuntimeError(nodes[1].Pos(), "too many arguments provided")
	}

	callbackBlock := runtime.NewBlock([]parser.Node{nodes[2]}, runtime.NewScope(block.Scope))

	for i, item := range l.List {
		if len(args) >= 1 {
			callbackBlock.Scope.SetSymbol(args[0], item)
		}

		if len(args) == 2 {
			callbackBlock.Scope.SetSymbol(args[1], runtime.NewNumberValueFromInt64(int64(i)))
		}

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

func stdMod(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromInt64(context.Args[0].NumberToInt64() % context.Args[1].NumberToInt64()), nil
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
		return nil, runtime.NewRuntimeError(context.Pos, "invalid argument(s), low can't be higher than high (%d > %d)", low, high)
	}

	l := runtime.NewListValue()

	for i := low; i < high; i++ {
		l.List = append(l.List, runtime.NewNumberValueFromInt64(i))
	}

	return l, nil
}

func stdPass(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue)

	if err != nil {
		return runtime.Nil, nil
	}

	return context.Args[0], nil
}

func stdLoad(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.StringValue)

	if err != nil {
		return nil, err
	}

	file, err := util.NewFile(context.Args[0].Str)

	if err != nil {
		return nil, err
	}

	l := lexer.NewLexer(lexer.NewSourceFromFile(file))
	util.ReportError(l.Lex(), false)

	p := parser.NewParser(l.Tokens)
	util.ReportError(p.Parse(), false)

	b := runtime.NewBlock(p.Nodes, context.Block.Scope)

	return b.Eval()
}

func stdSimpleMath(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	var callback func(float64) float64

	switch context.Name {
	case "sqrt":
		callback = math.Sqrt
	case "sin":
		callback = math.Sin
	case "cos":
		callback = math.Cos
	case "tan":
		callback = math.Tan
	case "ceil":
		callback = math.Ceil
	case "floor":
		callback = math.Floor
	case "abs":
		callback = math.Abs
	case "log":
		callback = math.Log
	case "log10":
		callback = math.Log10
	}

	return runtime.NewNumberValueFromFloat64(callback(context.Args[0].NumberToFloat64())), nil
}

func stdPow(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromFloat64(math.Pow(context.Args[0].NumberToFloat64(), context.Args[1].NumberToFloat64())), nil
}

func stdDeg2Rad(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromFloat64((context.Args[0].NumberToFloat64() * math.Pi) / 180), nil
}

func stdRad2Deg(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromFloat64((context.Args[0].NumberToFloat64() * 180) / math.Pi), nil
}
