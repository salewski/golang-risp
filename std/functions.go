package std

import "github.com/raoulvdberge/risp/runtime"

var Symbols = runtime.Symtab{
	"range":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdRange, "range"))),
	"substring": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdSubstring, "substring"))),
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
