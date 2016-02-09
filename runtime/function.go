package runtime

import "github.com/raoulvdberge/risp/lexer"

type BuiltinFunction func(*FunctionCallContext) (*Value, error)

type FunctionType int

const (
	Builtin FunctionType = iota
	Declared
)

type Function struct {
	Type         FunctionType
	Builtin      BuiltinFunction
	Declared     *Block
	DeclaredArgs []string
	Name         string
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
	case Declared:
		return f.Declared.Eval()
	}

	return nil, nil
}

func NewBuiltinFunction(function BuiltinFunction, name string) *Function {
	return &Function{Type: Builtin, Builtin: function, Name: name}
}

func NewDeclaredFunction(block *Block, name string, args []string) *Function {
	return &Function{Type: Declared, Declared: block, DeclaredArgs: args, Name: name}
}

type FunctionCallContext struct {
	Block *Block
	Args  []*Value
	Name  string
	Pos   *lexer.TokenPos
}
