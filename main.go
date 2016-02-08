package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/std"
	"github.com/raoulvdberge/risp/util"
	"os"
)

var (
	ast = flag.Bool("ast", false, "dump the abstract syntax tree")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) > 0 {
		runFile(flag.Arg(0))
	} else {
		runRepl()
	}
}

func runFile(name string) {
	file, err := util.NewFile(name)

	if err != nil {
		reportError(err, false)
	}

	l := lexer.NewLexerFromFile(file)
	reportError(l.Lex(), false)

	p := parser.NewParser(l.Tokens)
	reportError(p.Parse(), false)

	if *ast {
		fmt.Println(p.ToJson())
	} else {
		b := runtime.NewBlock(p.Nodes, runtime.NewScope(nil))
		b.Scope.Apply(std.Symbols)

		_, err := b.Eval()

		if err != nil {
			reportError(err, false)
		}
	}
}

func runRepl() {
	var tokens []*lexer.Token

	b := runtime.NewBlock(nil, runtime.NewScope(nil))
	b.Scope.Apply(std.Symbols)

	depth := 0

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		data, _ := reader.ReadString('\n')

		l := lexer.NewLexerFromString(data)

		err := l.Lex()

		if err != nil {
			reportError(err, true)
		} else {
			for _, t := range l.Tokens {
				tokens = append(tokens, t)

				depth += t.DepthModifier()
			}

			if depth <= 0 {
				p := parser.NewParser(tokens)

				err := p.Parse()

				tokens = nil
				depth = 0

				if err != nil {
					reportError(err, true)
				} else {
					b.Nodes = p.Nodes

					result, err := b.Eval()

					if err != nil {
						reportError(err, true)
					} else {
						fmt.Println(util.Yellow("===> " + result.String()))
					}
				}
			}
		}
	}
}

func reportError(err error, recoverable bool) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		if !recoverable {
			os.Exit(1)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: risp [options] [file ...]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
