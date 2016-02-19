package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/raoulvdberge/risp/builtin"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/list"
	"github.com/raoulvdberge/risp/math"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/repl"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/strings"
	"github.com/raoulvdberge/risp/util"
	"io/ioutil"
	"os"
)

var (
	runRepl = flag.Bool("repl", false, "runs the repl")
	ast     = flag.Bool("ast", false, "dumps the abstract syntax tree")
	debug   = flag.Bool("debug", false, "enabled debug mode")
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
	} else if *runRepl {
		s := repl.NewReplSession(apply(runtime.NewBlock(nil, runtime.NewScope(nil))))
		s.Run()
	} else {
		bytes, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			util.ReportError(err, false)
		}

		run(lexer.NewSourceFromString("<stdin>", string(bytes)))
	}
}

func apply(block *runtime.Block) *runtime.Block {
	block.Scope.ApplySymbols("", builtin.Symbols)
	block.Scope.ApplyMacros("", builtin.Macros)

	block.Scope.ApplySymbols("list", list.Symbols)
	block.Scope.ApplyMacros("list", list.Macros)

	block.Scope.ApplySymbols("string", strings.Symbols) // string is a type in Go so we have to keep using "strings" internally

	block.Scope.ApplySymbols("math", math.Symbols)

	return block
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
		b := apply(runtime.NewBlock(p.Nodes, runtime.NewScope(nil)))

		util.Timed("runtime", *debug, func() {
			_, err := b.Eval()

			if err != nil {
				util.ReportError(err, false)
			}
		})
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: risp [options] [file ...]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
