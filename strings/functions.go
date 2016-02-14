package strings

import (
	"github.com/raoulvdberge/risp/runtime"
	"strings"
)

var Symbols = runtime.Symtab{
	"substr":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsSubstr, "substr"))),
	"trim":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsTrim, "trim"))),
	"rune-at": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsRuneAt, "rune-at"))),
	"length":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsLength, "length"))),
}

func stringsSubstr(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.StringValue, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	source := context.Args[0].Str
	start := context.Args[1].NumberToInt64()
	end := context.Args[2].NumberToInt64()

	sourceLen := int64(len(source))

	if start < 0 || start > sourceLen || end < 0 || end > sourceLen {
		return nil, runtime.NewRuntimeError(context.Pos, "out of bounds (length is %d, trying to access %d:%d)", sourceLen, start, end)
	}

	result := source[start:end]

	return runtime.NewStringValue(result), nil
}

func stringsTrim(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue, runtime.StringValue); err != nil {
		return nil, err
	}

	source := context.Args[0].Str
	cutset := context.Args[1].Str
	return runtime.NewStringValue(strings.Trim(source, cutset)), nil
}

func stringsRuneAt(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue, runtime.NumberValue); err != nil {
		return nil, err
	}

	source := context.Args[0].Str
	index := context.Args[1].NumberToInt64()
	return runtime.NewStringValue(string([]rune(source)[index])), nil
}

func stringsLength(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue); err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromInt64(int64(len(context.Args[0].Str))), nil
}
