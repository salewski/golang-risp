package runtime

import (
	"github.com/raoulvdberge/risp/parser"
	"math/big"
)

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
	QuotedValue
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
	case QuotedValue:
		return "quoted"
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
	Quoted   parser.Node
}

func (v *Value) NumberToFloat64() float64 {
	f, _ := v.Number.Float64()

	return f
}

func (v *Value) NumberToInt64() int64 {
	return int64(v.NumberToFloat64())
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
	case FunctionValue:
		return v.Function.Name
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
	case QuotedValue:
		return v.Quoted.String()
	default:
		return "<" + v.Type.String() + ">"
	}
}

func (v *Value) Copy() *Value {
	if v.Type == NilValue {
		return Nil
	}

	other := &Value{Type: v.Type}

	switch v.Type {
	case StringValue:
		other.Str = v.Str
	case NumberValue:
		other.Number = v.Number
	case BooleanValue:
		other.Boolean = v.Boolean
	case KeywordValue:
		other.Keyword = v.Keyword
	case ListValue:
		for _, item := range v.List {
			other.List = append(other.List, item.Copy())
		}
	case FunctionValue:
		other.Function = v.Function.Copy()
	case QuotedValue:
		other.Quoted = v.Quoted
	}

	return other
}

func (v *Value) Equals(other *Value) bool {
	if v.Type == QuotedValue && other.Type == StringValue {
		return v.Quoted.String() == other.Str
	}

	if v.Type == StringValue && other.Type == QuotedValue {
		return other.Quoted.String() == v.Str
	}

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
	case QuotedValue:
		return v.Quoted.String() == other.Quoted.String()
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

func NewNumberValueFromFloat64(value float64) *Value {
	number := big.NewRat(0, 1)
	number.SetFloat64(value)

	return &Value{Type: NumberValue, Number: number}
}

func NewNumberValueFromInt64(value int64) *Value {
	number := big.NewRat(0, 1)
	number.SetInt64(value)

	return &Value{Type: NumberValue, Number: number}
}

func NewFunctionValue(value *Function) *Value {
	return &Value{Type: FunctionValue, Function: value}
}

func NewListValue() *Value {
	return &Value{Type: ListValue}
}

func NewQuotedValue(node parser.Node) *Value {
	return &Value{Type: QuotedValue, Quoted: node}
}

func BooleanValueFor(value bool) *Value {
	if value == true {
		return True
	} else {
		return False
	}
}
