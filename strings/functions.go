package strings

import "github.com/raoulvdberge/risp/runtime"

var Symbols = runtime.Symtab{
	"substr": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsSubstr, "substr"))),
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
