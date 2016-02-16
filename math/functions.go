package math

import (
	"github.com/raoulvdberge/risp/runtime"
	"math"
)

var Symbols = runtime.Symtab{
	"mod":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathMod, "mod"))),
	"sqrt":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "sqrt"))),
	"sin":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "sin"))),
	"cos":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "cos"))),
	"tan":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "tan"))),
	"ceil":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "ceil"))),
	"floor":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "floor"))),
	"abs":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "abs"))),
	"log":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "log"))),
	"log10":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathSimpleMath, "log10"))),
	"pow":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathPow, "pow"))),
	"deg2rad": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathDeg2Rad, "deg2rad"))),
	"rad2deg": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(mathRad2Deg, "rad2deg"))),
	"pi":      runtime.NewSymbol(runtime.NewNumberValueFromFloat64(math.Pi)),
	"e":       runtime.NewSymbol(runtime.NewNumberValueFromFloat64(math.E)),
}

func mathMod(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue); err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromInt64(context.Args[0].NumberToInt64() % context.Args[1].NumberToInt64()), nil
}

func mathSimpleMath(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.NumberValue); err != nil {
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

func mathPow(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue); err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromFloat64(math.Pow(context.Args[0].NumberToFloat64(), context.Args[1].NumberToFloat64())), nil
}

func mathDeg2Rad(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.NumberValue); err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromFloat64((context.Args[0].NumberToFloat64() * math.Pi) / 180), nil
}

func mathRad2Deg(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.NumberValue); err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromFloat64((context.Args[0].NumberToFloat64() * 180) / math.Pi), nil
}
