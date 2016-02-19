package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/peterh/liner"
	"github.com/raoulvdberge/risp/builtin"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/list"
	"github.com/raoulvdberge/risp/math"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/strings"
	"github.com/raoulvdberge/risp/util"
	"io/ioutil"
	"os"
	"strconv"
	stringsPkg "strings"
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

func apply(scope *runtime.Scope) {
	scope.ApplySymbols("", builtin.Symbols)
	scope.ApplyMacros("", builtin.Macros)

	scope.ApplySymbols("list", list.Symbols)
	scope.ApplyMacros("list", list.Macros)

	scope.ApplySymbols("string", strings.Symbols) // string is a type in Go so we have to keep using "strings" internally

	scope.ApplySymbols("math", math.Symbols)
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

		apply(b.Scope)

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

	apply(b.Scope)

	line.SetCompleter(func(line string) (c []string) {
		ident := ""
		identEnd := -1

		for i := len(line) - 1; i >= 0; i-- {
			identEnd = i

			part := rune(line[i])

			if lexer.IsIdentifierPart(part) {
				ident += string(part)
			} else if lexer.IsIdentifierStart(part) {
				ident += string(part)

				break
			} else {
				break
			}
		}

		ident = util.ReverseString(ident)

		if ident != "" && lexer.IsIdentifierStart(rune(ident[0])) && identEnd != -1 {
			for key, _ := range b.Scope.Symbols {
				if stringsPkg.HasPrefix(key, ident) {
					c = append(c, line[0:identEnd+1]+key)
				}
			}
		}

		return
	})

	depth := 0

	for {
		prompt := "> "

		if depth > 0 {
			prompt += "(" + strconv.Itoa(depth) + ") "
		}

		data, err := line.Prompt(prompt)

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

							b.Scope.SetSymbolLocally("_", runtime.NewSymbol(result))
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
