package list

import "github.com/raoulvdberge/risp/runtime"

var Symbols = runtime.Symtab{
	"seq":          runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listSeq, "seq"))),
	"contains":     runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listContains, "contains"))),
	"contains-key": runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listContainsKey, "contains-key"))),
	"push":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listPush, "push"))),
	"push-left":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listPushLeft, "push-left"))),
	"size":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listSize, "size"))),
	"get":          runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listGet, "get"))),
	"get-key":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listGetKey, "get-key"))),
	"set":          runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listSet, "set"))),
	"set-key":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listSetKey, "set-key"))),
	"drop":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listDrop, "drop"))),
	"drop-left":    runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listDropLeft, "drop-left"))),
	"join":         runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listJoin, "join"))),
	"range":        runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listRange, "range"))),
	"reverse":      runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listReverse, "reverse"))),
	"remove":       runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listRemove, "remove"))),
	"remove-key":   runtime.NewSymbol(runtime.NewFunctionValue(runtime.NewBuiltinFunction(listRemoveKey, "remove-key"))),
}

func listSeq(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue); err != nil {
		return nil, err
	}

	low := context.Args[0].NumberToInt64()
	high := context.Args[1].NumberToInt64()

	if low > high {
		return nil, runtime.NewRuntimeError(context.Pos, "invalid argument(s), low can't be higher than high (%d > %d)", low, high)
	}

	l := runtime.NewListValue()

	for i := low; i <= high; i++ {
		l.List = append(l.List, runtime.NewNumberValueFromInt64(i))
	}

	return l, nil
}

func listContains(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.AnyValue); err != nil {
		return nil, err
	}

	for _, item := range context.Args[0].List {
		if item.Equals(context.Args[1]) {
			return runtime.True, nil
		}
	}

	return runtime.False, nil
}

func listContainsKey(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.KeywordValue); err != nil {
		return nil, err
	}

	l := context.Args[0].List

	for i := 0; i < len(l); i++ {
		if l[i].Type == runtime.KeywordValue && l[i].Keyword == context.Args[1].Keyword && i+1 < len(l) {
			return runtime.True, nil
		}
	}

	return runtime.False, nil
}

func listPush(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.AnyValue); err != nil {
		return nil, err
	}

	context.Args[0].List = append(context.Args[0].List, context.Args[1])

	return context.Args[0], nil
}

func listPushLeft(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.AnyValue); err != nil {
		return nil, err
	}

	context.Args[0].List = append([]*runtime.Value{context.Args[1]}, context.Args[0].List...)

	return context.Args[0], nil
}

func listSize(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue); err != nil {
		return nil, err
	}

	return runtime.NewNumberValueFromInt64(int64(len(context.Args[0].List))), nil
}

func listGet(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.NumberValue); err != nil {
		return nil, err
	}

	size := int64(len(context.Args[0].List))
	index := context.Args[1].NumberToInt64()

	if index < 0 || index > size-1 {
		return nil, runtime.NewRuntimeError(context.Pos, "index %d out of bounds (list size is %d)", index, size)
	}

	return context.Args[0].List[index], nil
}

func listGetKey(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.KeywordValue); err != nil {
		return nil, err
	}

	l := context.Args[0]

	for i := 0; i < len(l.List); i++ {
		if l.List[i].Type == runtime.KeywordValue && l.List[i].Keyword == context.Args[1].Keyword && i+1 < len(l.List) {
			return l.List[i+1], nil
		}
	}

	return runtime.Nil, nil
}

func listSet(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.NumberValue, runtime.AnyValue); err != nil {
		return nil, err
	}

	size := int64(len(context.Args[0].List))
	index := context.Args[1].NumberToInt64()

	if index < 0 || index > size-1 {
		return nil, runtime.NewRuntimeError(context.Pos, "index %d out of bounds (list size is %d)", index, size)
	}

	context.Args[0].List[index] = context.Args[2]

	return context.Args[0], nil
}

func listSetKey(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.KeywordValue, runtime.AnyValue); err != nil {
		return nil, err
	}

	l := context.Args[0]

	for i := 0; i < len(l.List); i++ {
		if l.List[i].Type == runtime.KeywordValue && l.List[i].Keyword == context.Args[1].Keyword && i+1 < len(l.List) {
			l.List[i+1] = context.Args[2]

			return l, nil
		}
	}

	l.List = append(l.List, context.Args[1])
	l.List = append(l.List, context.Args[2])

	return l, nil
}

func listDrop(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue); err != nil {
		return nil, err
	}

	l := context.Args[0]

	if len(l.List) == 0 {
		return nil, runtime.NewRuntimeError(context.Pos, "empty list")
	}

	l.List = l.List[0 : len(l.List)-1]

	return context.Args[0], nil
}

func listDropLeft(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue); err != nil {
		return nil, err
	}

	l := context.Args[0]

	if len(l.List) == 0 {
		return nil, runtime.NewRuntimeError(context.Pos, "empty list")
	}

	l.List = l.List[1:]

	return context.Args[0], nil
}

func listJoin(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.ListValue); err != nil {
		return nil, err
	}

	for _, item := range context.Args[1].List {
		context.Args[0].List = append(context.Args[0].List, item)
	}

	return context.Args[0], nil
}

func listRange(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	var begin, end int64

	err := runtime.ValidateArguments(context, runtime.ListValue, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		optionalErr := runtime.ValidateArguments(context, runtime.ListValue, runtime.NumberValue)

		if optionalErr == nil {
			begin = context.Args[1].NumberToInt64()
			end = int64(len(context.Args[0].List)) - 1
		} else {
			return nil, optionalErr
		}
	} else {
		begin = context.Args[1].NumberToInt64()
		end = context.Args[2].NumberToInt64()
	}

	length := int64(len(context.Args[0].List))

	if begin < 0 || begin > length-1 || begin > end || end < 0 || end > length-1 {
		return nil, runtime.NewRuntimeError(context.Pos, "invalid bounds %d and %d (list length is %d)", begin, end, length)
	}

	newList := runtime.NewListValue()

	for _, item := range context.Args[0].List[begin : end+1] {
		newList.List = append(newList.List, item)
	}

	if len(newList.List) == 1 {
		return newList.List[0], nil
	} else {
		return newList, nil
	}
}

func listReverse(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue); err != nil {
		return nil, err
	}

	list := context.Args[0].List
	newList := runtime.NewListValue()

	for i := len(list) - 1; i >= 0; i-- {
		newList.List = append(newList.List, list[i])
	}

	return newList, nil
}

func listRemove(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.NumberValue); err != nil {
		return nil, err
	}

	list := context.Args[0].List
	size := int64(len(list))
	index := context.Args[1].NumberToInt64()

	if index < 0 || index > size-1 {
		return nil, runtime.NewRuntimeError(context.Pos, "index %d out of bounds (list size is %d)", index, size)
	}

	newList := runtime.NewListValue()

	for i, item := range list {
		if int64(i) != index {
			newList.List = append(newList.List, item)
		}
	}

	return newList, nil
}

func listRemoveKey(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	if err := runtime.ValidateArguments(context, runtime.ListValue, runtime.KeywordValue); err != nil {
		return nil, err
	}

	l := context.Args[0]

	for i := 0; i < len(l.List); i++ {
		if l.List[i].Type == runtime.KeywordValue && l.List[i].Keyword == context.Args[1].Keyword && i+1 < len(l.List) {
			l.List = append(l.List[:i], l.List[i+1:]...)
			l.List = append(l.List[:i], l.List[i+1:]...)
		}
	}

	return l, nil
}
