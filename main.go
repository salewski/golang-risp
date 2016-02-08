package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/std"
	"github.com/raoulvdberge/risp/util"
	"io/ioutil"
	"os"
)

var (
	ast       = flag.Bool("ast", false, "dump the abstract syntax tree")
	astIndent = flag.String("ast-indent", "", "the indendentation character for the abstract syntax tree output")
	fromStdin = flag.Bool("from-stdin", false, "read input from stdin")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if *fromStdin {
		bytes, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			reportError(err, false)
		}

		run(lexer.NewSourceFromString("<stdin>", string(bytes)))
	} else if len(flag.Args()) > 0 {
		file, err := util.NewFile(flag.Arg(0))

		if err != nil {
			reportError(err, false)
		}

		run(lexer.NewSourceFromFile(file))
	} else {
		runRepl()
	}
}

func run(source lexer.Source) {
	l := lexer.NewLexer(source)
	reportError(l.Lex(), false)

	p := parser.NewParser(l.Tokens)
	reportError(p.Parse(), false)

	if *ast {
		var bytes []byte

		if *astIndent == "" {
			bytes, _ = json.Marshal(p)
		} else {
			bytes, _ = json.MarshalIndent(p, "", *astIndent)
		}

		fmt.Println(string(bytes))
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

		l := lexer.NewLexer(lexer.NewSourceFromString("<repl>", data))

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
