package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/peterh/liner"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/std"
	"github.com/raoulvdberge/risp/util"
	"io/ioutil"
	"os"
	"strings"
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
		b.Scope.ApplySymbols(runtime.Primitives)

		b.Scope.ApplyMacros(std.Macros)

		util.Timed("runtime", *debug, func() {
			_, err := b.Eval()

			if err != nil {
				util.ReportError(err, false)
			}
		})
	}
}

func runRepl() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	var tokens []*lexer.Token

	b := runtime.NewBlock(nil, runtime.NewScope(nil))

	b.Scope.ApplySymbols(std.Symbols)
	b.Scope.ApplySymbols(runtime.Primitives)

	b.Scope.ApplyMacros(std.Macros)

	line.SetCompleter(func(line string) (c []string) {
		l := lexer.NewLexer(lexer.NewSourceFromString("", line))
		l.Formatting = true
		l.Lex()

		if len(l.Tokens) > 0 {
			prev := ""

			for i := 0; i < len(l.Tokens)-1; i++ {
				prev += l.Tokens[i].Data
			}

			last := l.Tokens[len(l.Tokens)-1]

			if last.Type == lexer.Identifier {
				for key, _ := range b.Scope.Symbols {
					if strings.HasPrefix(key, last.Data) {
						c = append(c, prev+key)
					}
				}

				for key, _ := range b.Scope.Macros {
					if strings.HasPrefix(key, last.Data) {
						c = append(c, prev+key)
					}
				}
			}
		}

		return
	})

	depth := 0

	for {
		data, err := line.Prompt("> ")

		if err != nil {
			if err == liner.ErrPromptAborted {
				return
			}

			util.ReportError(err, false)
		}

		line.AppendHistory(data)

		l := lexer.NewLexer(lexer.NewSourceFromString("<repl>", data))

		err = l.Lex()

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
							data := util.Yellow("===> " + result.String())

							if result.Type != runtime.NilValue {
								data += " " + util.Yellow("("+result.Type.String()+")")
							}

							fmt.Println(data)

							b.Scope.SetSymbol("_", runtime.NewSymbol(result))
						}
					}
				}
			} else {
				fmt.Printf("%s===> %d missing closing parenthesis%s\n", util.TEXT_GREEN, depth, util.TEXT_RESET)
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
