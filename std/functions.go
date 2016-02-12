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
	"print":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrint, "print"))),
	"println":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrintln, "println"))),
	"list":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdList, "list"))),
	"+":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "+"))),
	"-":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "-"))),
	"*":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "*"))),
	"/":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "/"))),
	"%":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPow, "%"))),
	"=":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdEquals, "="))),
	"!=":        runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdNotEquals, "!="))),
	">":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, ">"))),
	">=":        runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, ">="))),
	"<":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, "<"))),
	"<=":        runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, "<="))),
	"cat":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdCat, "cat"))),
	"and":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdAnd, "and"))),
	"or":        runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdOr, "or"))),
	"not":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdNot, "not"))),
	"call":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdCall, "call"))),
	"range":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdRange, "range"))),
	"pass":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPass, "pass"))),
	"load":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdLoad, "load"))),
	"sqrt":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "sqrt"))),
	"sin":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "sin"))),
	"cos":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "cos"))),
	"tan":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "tan"))),
	"ceil":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "ceil"))),
	"floor":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "floor"))),
	"abs":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "abs"))),
	"log":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "log"))),
	"log10":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSimpleMath, "log10"))),
	"pow":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPow, "pow"))),
	"deg2rad":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdDeg2Rad, "deg2rad"))),
	"rad2deg":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdRad2Deg, "rad2deg"))),
	"pi":        runtime.NewSymbol(runtime.NewNumberValueFromFloat64(math.Pi)),
	"e":         runtime.NewSymbol(runtime.NewNumberValueFromFloat64(math.E)),
	"substring": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSubstring, "substring"))),
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
	if start < 0 || start > sourceLen || end < 0 || end > sourceLen {
		return nil, runtime.NewRuntimeError(context.Pos, "attempting to slice string out of bounds")
	}

	result := source[start:end]
	return runtime.NewStringValue(result), nil
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
