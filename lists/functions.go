package lists

import "github.com/raoulvdberge/risp/runtime"

var Symbols = runtime.Symtab{
	"range": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listsRange, "range"))),
	"push":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listsPush, "push"))),
}

func listsRange(context *runtime.FunctionCallContext) (*runtime.Value, error) {
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

func listsPush(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.ListValue, runtime.AnyValue)

	if err != nil {
		return nil, err
	}

	context.Args[0].List = append(context.Args[0].List, context.Args[1])

	return context.Args[0], nil
}
