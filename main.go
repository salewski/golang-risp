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
	repl  = flag.Bool("repl", false, "runs the repl")
	ast   = flag.Bool("ast", false, "dumps the abstract syntax tree")
	debug = flag.Bool("debug", false, "enabled debug mode")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) > 0 {
		var file *util.File

		util.Timed("file reading", *debug, func() {
			f, err := util.NewFile(flag.Arg(0))

			if err != nil {
				util.ReportError(err, false)
			}

			file = f
		})

		run(lexer.NewSourceFromFile(file))
	} else if *repl {
		runRepl()
	} else {
		bytes, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			util.ReportError(err, false)
		}

		run(lexer.NewSourceFromString("<stdin>", string(bytes)))
	}
}

func run(source lexer.Source) {
	l := lexer.NewLexer(source)
	util.Timed("lexing", *debug, func() {
		util.ReportError(l.Lex(), false)
	})

	p := parser.NewParser(l.Tokens)
	util.Timed("parsing", *debug, func() {
		util.ReportError(p.Parse(), false)
	})

	if *ast {
		bytes, _ := json.MarshalIndent(p, "", "    ")

		fmt.Println(string(bytes))
	} else {
		b := runtime.NewBlock(p.Nodes, runtime.NewScope(nil))
		b.Scope.ApplySymbols(std.Symbols)
		b.Scope.ApplyMacros(std.Macros)
		b.Scope.ApplySymbols(runtime.Primitives)

		util.Timed("runtime", *debug, func() {
			_, err := b.Eval()

			if err != nil {
				util.ReportError(err, false)
			}
		})
	}
}

func runRepl() {
	var tokens []*lexer.Token

	b := runtime.NewBlock(nil, runtime.NewScope(nil))
	b.Scope.ApplySymbols(std.Symbols)
	b.Scope.ApplyMacros(std.Macros)
	b.Scope.ApplySymbols(runtime.Primitives)

	depth := 0

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")

		data, _ := reader.ReadString('\n')

		l := lexer.NewLexer(lexer.NewSourceFromString("<repl>", data))

		err := l.Lex()

		if err != nil {
			util.ReportError(err, true)
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
					util.ReportError(err, true)
				} else {
					for _, node := range p.Nodes {
						result, err := b.EvalNode(node)

						if err != nil {
							util.ReportError(err, true)
						} else {
							fmt.Println(util.Yellow("===> " + result.String() + " (" + result.Type.String() + ")"))

							b.Scope.SetSymbol("_", result)
						}
					}
				}
			}
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: risp [options] [file ...]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
