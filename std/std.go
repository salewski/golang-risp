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
	"substring": runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSubstring, "substring")),
}

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

func stdSubstring(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.StringValue, runtime.NumberValue, runtime.NumberValue)
	if err != nil {
		return nil, err
	}

	source := context.Args[0].Str
	start := context.Args[1].NumberToInt64()
	end := context.Args[2].NumberToInt64()

	sourceLen := int64(len(source))
	// verify the range isn't shit
	if (start < 0 || start > sourceLen || end < 0 || end > sourceLen) {
		return nil, runtime.NewRuntimeError(context.Pos, "attempting to slice string out of bounds")
	}

	result := source[start:end]
	return runtime.NewStringValue(result), nil
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

	context.Block.Scope.SetSymbol(name, runtime.NewFunctionValue(function))

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

	context.Block.Scope.SetSymbol(name, value)

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
