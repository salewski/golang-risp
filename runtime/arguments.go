package runtime

func ValidateArguments(context *FunctionCallContext, types ...ValueType) error {
	if len(types) != len(context.Args) {
		return NewRuntimeError(context.Pos, "%s: expected %d arguments, got %d", context.Name, len(types), len(context.Args))
	}

	for i, arg := range context.Args {
		if types[i] != AnyValue {
			if arg.Type != types[i] {
				return NewRuntimeError(context.Pos, "%s: argument %d should be of type %s, got %s", context.Name, i+1, types[i], arg.Type)
			}
		}
	}

	return nil
}
