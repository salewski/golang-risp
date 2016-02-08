package runtime

import "math/big"

var (
	Nil   = &Value{Type: NilValue}
	True  = &Value{Type: BooleanValue, Boolean: true}
	False = &Value{Type: BooleanValue, Boolean: false}
)

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
		return ":" + v.Keyword
	case ListValue:
		s := "("

		for i, item := range v.List {
			s += item.String()

			if i != len(v.List)-1 {
				s += " "
			}
		}

		s += ")"

		return s
	case NilValue:
		return "nil"
	default:
		return "<" + v.Type.String() + ">"
	}
}

func (v *Value) Equals(other *Value) bool {
	if v.Type != other.Type {
		return false
	}

	switch v.Type {
	case StringValue:
		return v.Str == other.Str
	case NumberValue:
		return v.Number.Cmp(other.Number) == 0
	case BooleanValue:
		return v.Boolean == other.Boolean
	case KeywordValue:
		return v.Keyword == other.Keyword
	case ListValue:
		if len(v.List) != len(other.List) {
			return false
		}

		for i, v := range v.List {
			if !v.Equals(other.List[i]) {
				return false
			}
		}

		return true
	case FunctionValue:
		return &v.Function == &other.Function
	case NilValue:
		return true
	default:
		return false
	}
}

func NewStringValue(value string) *Value {
	return &Value{Type: StringValue, Str: value}
}

func NewKeywordValue(value string) *Value {
	return &Value{Type: KeywordValue, Keyword: value}
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

func NewListValue() *Value {
	return &Value{Type: ListValue}
}
