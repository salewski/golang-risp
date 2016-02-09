package runtime

import "github.com/raoulvdberge/risp/lexer"

type BuiltinFunction func(*FunctionCallContext) (*Value, error)

type FunctionType int

const (
	Builtin FunctionType = iota
	Declared
	Lambda
)

type Function struct {
	Type FunctionType
	Name string
	// for builtin functions
	Builtin BuiltinFunction
	// for lambdas and declared functions
	Callback *Block
	Args     []string
}

func (f *Function) Call(block *Block, args []*Value, pos *lexer.TokenPos) (*Value, error) {
	switch f.Type {
	case Builtin:
		return f.Builtin(&FunctionCallContext{
			Block: block,
			Args:  args,
			Name:  f.Name,
			Pos:   pos,
		})
	case Declared, Lambda:
		if len(args) != len(f.Args) {
			return nil, NewRuntimeError(pos, "'%s' expected %d arguments, got %d", f.Name, len(f.Args), len(args))
		}

		for i, argName := range f.Args {
			f.Callback.Scope.SetSymbol(argName, args[i])
		}

		return f.Callback.Eval()
	}

	return nil, nil
}

func NewBuiltinFunction(function BuiltinFunction, name string) *Function {
	return &Function{Type: Builtin, Builtin: function, Name: name}
}

func NewDeclaredFunction(block *Block, name string, args []string) *Function {
	return &Function{Type: Declared, Callback: block, Args: args, Name: name}
}

func NewLambdaFunction(block *Block, args []string) *Function {
	return &Function{Type: Lambda, Callback: block, Args: args, Name: "<lambda>"}
}

type FunctionCallContext struct {
	Block *Block
	Args  []*Value
	Name  string
	Pos   *lexer.TokenPos
}
