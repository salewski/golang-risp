package builtin

import (
	"fmt"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/util"
	"math/big"
)

var Symbols = runtime.Symtab{
	"t":       runtime.NewSymbol(runtime.True),
	"f":       runtime.NewSymbol(runtime.False),
	"nil":     runtime.NewSymbol(runtime.Nil),
	"print":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinPrint, "print"))),
	"println": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinPrintln, "println"))),
	"list":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinList, "list"))),
	"+":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMath, "+"))),
	"-":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMath, "-"))),
	"*":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMath, "*"))),
	"/":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMath, "/"))),
	"=":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinEquals, "="))),
	"!=":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinNotEquals, "!="))),
	">":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMathCmp, ">"))),
	">=":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMathCmp, ">="))),
	"<":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMathCmp, "<"))),
	"<=":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinMathCmp, "<="))),
	"and":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinAnd, "and"))),
	"or":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinOr, "or"))),
	"not":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinNot, "not"))),
	"call":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinCall, "call"))),
	"pass":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinPass, "pass"))),
	"load":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinLoad, "load"))),
	"cat":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinCat, "cat"))),
	"assert":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(builtinAssert, "assert"))),
}

func builtinPrint(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Print(arg)
	}

	return runtime.Nil, nil
}

func builtinPrintln(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Println(arg)
	}

	return runtime.Nil, nil
}

func builtinList(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	l := runtime.NewListValue()

	for _, arg := range context.Args {
		l.List = append(l.List, arg)
	}

	return l, nil
}

func builtinMath(context *runtime.FunctionCallContext) (*runtime.Value, error) {
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

func builtinMathCmp(context *runtime.FunctionCallContext) (*runtime.Value, error) {
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

func builtinEquals(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue, runtime.AnyValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(context.Args[0].Equals(context.Args[1])), nil
}

func builtinNotEquals(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue, runtime.AnyValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(!context.Args[0].Equals(context.Args[1])), nil
}

func builtinAnd(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(context.Args[0].Boolean && context.Args[1].Boolean), nil
}

func builtinOr(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(context.Args[0].Boolean || context.Args[1].Boolean), nil
}

func builtinNot(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(!context.Args[0].Boolean), nil
}

func builtinCall(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if len(context.Args) < 1 {
		return nil, runtime.NewRuntimeError(context.Pos, "expected a function")
	}

	function := context.Args[0].Function

	return function.Call(context.Block, context.Args[1:], context.Pos)
}

func builtinPass(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue)

	if err != nil {
		return runtime.Nil, nil
	}

	return context.Args[0], nil
}

func builtinLoad(context *runtime.FunctionCallContext) (*runtime.Value, error) {
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

	b := runtime.NewBlock(p.Nodes, runtime.NewScope(context.Block.Scope))

	result, err := b.Eval()

	if err != nil {
		return nil, err
	}

	for key, value := range b.Scope.Symbols {
		if value.Exported {
			if value.Value.Type == runtime.FunctionValue {
				value.Value.Function.CustomScope = b.Scope
			}

			context.Block.Scope.SetSymbol(b.SymbolName(key), value)
		}
	}

	return result, nil
}

func builtinCat(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	s := ""

	for _, arg := range context.Args {
		s += arg.String()
	}

	return runtime.NewStringValue(s), nil
}

func builtinAssert(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	optionalError := runtime.ValidateArguments(context, runtime.BooleanValue)

	if optionalError != nil {
		err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.StringValue)

		if err != nil {
			return nil, err
		}
	}

	message := "assertion failed"

	if len(context.Args) == 2 {
		message += ": " + context.Args[1].Str
	}

	if !context.Args[0].Boolean {
		return nil, runtime.NewRuntimeError(context.Pos, message)
	}

	return runtime.Nil, nil
}
