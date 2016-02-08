package std

import (
	"fmt"
	"github.com/raoulvdberge/risp/runtime"
	"math/big"
)

func Apply(block *runtime.Block) {
	block.Scope.Symbols["t"] = runtime.NewBooleanValue(true)
	block.Scope.Symbols["f"] = runtime.NewBooleanValue(false)
	block.Scope.Symbols["nil"] = runtime.Nil
	block.Scope.Symbols["print"] = runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrint, "print"))
	block.Scope.Symbols["list"] = runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdList, "list"))
	block.Scope.Symbols["+"] = runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdAdd, "+"))
}

func stdPrint(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Print(arg)
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

func stdAdd(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromRat(big.NewRat(0, 1).Add(context.Args[0].Number, context.Args[1].Number)), nil
}
