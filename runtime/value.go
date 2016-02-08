package runtime

import "math/big"

var Nil = &Value{Type: NilValue}

type ValueType int

const (
	StringValue ValueType = iota
	NumberValue
	BooleanValue
	KeywordValue
	ListValue
	FunctionValue
	NilValue
	AnyValue // used in arguments.go, to validate *any* argument
)

func (t ValueType) String() string {
	switch t {
	case StringValue:
		return "string"
	case NumberValue:
		return "number"
	case BooleanValue:
		return "boolean"
	case KeywordValue:
		return "keyword"
	case ListValue:
		return "list"
	case FunctionValue:
		return "function"
	case NilValue:
		return "nil"
	default:
		return "?"
	}
}

type Value struct {
	Type     ValueType
	Str      string
	Number   *big.Rat
	Boolean  bool
	Keyword  string
	List     []*Value
	Function *Function
}

func (v *Value) String() string {
	switch v.Type {
	case StringValue:
		return v.Str
	case NumberValue:
		prec := 4

		if v.Number.Denom().Cmp(big.NewInt(1)) == 0 {
			prec = 0
		}

		return v.Number.FloatString(prec)
	case BooleanValue:
		if v.Boolean {
			return "t"
		} else {
			return "f"
		}
	case KeywordValue:
		return v.Keyword
	case ListValue:
		// @TODO
		return "<list>"
	case FunctionValue:
		return "<function>"
	case NilValue:
		return "nil"
	default:
		return "<?>"
	}
}

func NewStringValue(value string) *Value {
	return &Value{Type: StringValue, Str: value}
}

func NewKeywordValue(value string) *Value {
	return &Value{Type: KeywordValue, Str: value}
}

func NewNumberValueFromString(value string) *Value {
	number := big.NewRat(0, 1)
	number.SetString(value)

	return NewNumberValueFromRat(number)
}

func NewNumberValueFromRat(value *big.Rat) *Value {
	return &Value{Type: NumberValue, Number: value}
}

func NewFunctionValue(value *Function) *Value {
	return &Value{Type: FunctionValue, Function: value}
}

func NewBooleanValue(value bool) *Value {
	return &Value{Type: BooleanValue, Boolean: value}
}
