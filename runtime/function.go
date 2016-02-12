package runtime

import (
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
)

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
	Nodes []parser.Node
	Args  []string
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
		functionBlock := NewBlock(f.Nodes, NewScope(block.Scope))

		if len(args) != len(f.Args) {
			return nil, NewRuntimeError(pos, "'%s' expected %d arguments, got %d", f.Name, len(f.Args), len(args))
		}

		for i, argName := range f.Args {
			functionBlock.Scope.SetSymbol(argName, args[i])
		}

		return functionBlock.Eval()
	}

	return nil, nil
}

func NewBuiltinFunction(function BuiltinFunction, name string) *Function {
	return &Function{Type: Builtin, Builtin: function, Name: name}
}

func NewDeclaredFunction(nodes []parser.Node, name string, args []string) *Function {
	return &Function{Type: Declared, Nodes: nodes, Args: args, Name: name}
}

func NewLambdaFunction(nodes []parser.Node, args []string) *Function {
	return &Function{Type: Lambda, Nodes: nodes, Args: args, Name: "<lambda>"}
}

type FunctionCallContext struct {
	Block *Block
	Args  []*Value
	Name  string
	Pos   *lexer.TokenPos
}
