package strings

import (
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/util"
	"strings"
	"unicode"
)

var Symbols = runtime.Symtab{
	"range":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsRange, "range"))),
	"trim":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsTrim, "trim"))),
	"rune-at":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsRuneAt, "rune-at"))),
	"length":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsLength, "length"))),
	"format":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsFormat, "format"))),
	"split":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsSplit, "split"))),
	"replace":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsReplace, "replace"))),
	"reverse":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsReverse, "reverse"))),
	"contains":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsContains, "contains"))),
	"lower":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsLower, "lower"))),
	"upper":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsUpper, "upper"))),
	"is-digit":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsCharacterCheck, "is-digit"))),
	"is-letter": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsCharacterCheck, "is-letter"))),
	"is-lower":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsCharacterCheck, "is-lower"))),
	"is-upper":  runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(stringsCharacterCheck, "is-upper"))),
}

func stringsRange(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue, runtime.NumberValue, runtime.NumberValue); err != nil {
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

func stringsFormat(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if len(context.Args) < 1 {
		return nil, runtime.NewRuntimeError(context.Pos, "missing format specifier")
	}

	if context.Args[0].Type != runtime.StringValue {
		return nil, runtime.NewRuntimeError(context.Pos, "format specifier should be a string")
	}

	format := context.Args[0].Str

	args := strings.Count(format, "~")

	if len(context.Args)-1 != args {
		return nil, runtime.NewRuntimeError(context.Pos, "format specifier expected %d arguments, got %d", args, len(context.Args)-1)
	}

	for _, item := range context.Args[1:] {
		format = strings.Replace(format, "~", item.String(), 1)
	}

	return runtime.NewStringValue(format), nil
}

func stringsSplit(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue, runtime.StringValue); err != nil {
		return nil, err
	}

	parts := runtime.NewListValue()

	for _, item := range strings.Split(context.Args[0].Str, context.Args[1].Str) {
		parts.List = append(parts.List, runtime.NewStringValue(item))
	}

	return parts, nil
}

func stringsReplace(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	n := -1

	err := runtime.ValidateArguments(context, runtime.StringValue, runtime.StringValue, runtime.StringValue, runtime.NumberValue)

	if err != nil {
		optionalErr := runtime.ValidateArguments(context, runtime.StringValue, runtime.StringValue, runtime.StringValue)

		if optionalErr != nil {
			return nil, err
		}
	} else {
		n = int(context.Args[3].NumberToInt64())
	}

	source := context.Args[0].Str
	search := context.Args[1].Str
	replace := context.Args[2].Str

	return runtime.NewStringValue(strings.Replace(source, search, replace, n)), nil
}

func stringsReverse(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue); err != nil {
		return nil, err
	}

	return runtime.NewStringValue(util.ReverseString(context.Args[0].Str)), nil
}

func stringsContains(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue, runtime.StringValue); err != nil {
		return nil, err
	}

	return runtime.BooleanValueFor(strings.Contains(context.Args[0].Str, context.Args[1].Str)), nil
}

func stringsLower(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue); err != nil {
		return nil, err
	}

	return runtime.NewStringValue(strings.ToLower(context.Args[0].Str)), nil
}

func stringsUpper(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue); err != nil {
		return nil, err
	}

	return runtime.NewStringValue(strings.ToUpper(context.Args[0].Str)), nil
}

func stringsCharacterCheck(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.StringValue); err != nil {
		return nil, err
	}

	s := context.Args[0].Str

	if len(s) != 1 {
		return nil, runtime.NewRuntimeError(context.Pos, "expected string that has 1 character, got %d character(s)", len(s))
	}

	var callback func(rune) bool

	switch context.Name {
	case "is-digit":
		callback = unicode.IsDigit
	case "is-letter":
		callback = unicode.IsLetter
	case "is-lower":
		callback = unicode.IsLower
	case "is-upper":
		callback = unicode.IsUpper
	}

	return runtime.BooleanValueFor(callback(rune(s[0]))), nil
}
