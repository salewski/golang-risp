package repl

import (
	"fmt"
	"github.com/peterh/liner"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/util"
	"strconv"
)

type ReplSession struct {
	Block  *runtime.Block
	tokens []*lexer.Token
	depth  int
}

func NewReplSession(block *runtime.Block) *ReplSession {
	return &ReplSession{
		Block: block,
		depth: 0,
	}
}

func (s *ReplSession) Run() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(s.completer)

	for {
		prompt := "> "

		if s.depth > 0 {
			prompt += "(" + strconv.Itoa(s.depth) + ") "
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
				s.tokens = append(s.tokens, t)
				s.depth += t.DepthModifier()
			}

			if s.depth <= 0 {
				p := parser.NewParser(s.tokens)

				err := p.Parse()

				s.tokens = nil
				s.depth = 0

				if err != nil {
					util.ReportError(err, true)
				} else {
					for _, node := range p.Nodes {
						result, err := s.Block.EvalNode(node)

						if err != nil {
							util.ReportError(err, true)
						} else {
							data := util.Yellow("===> " + result.String())

							if result.Type != runtime.NilValue {
								data += " " + util.Yellow("("+result.Type.String()+")")
							}

							fmt.Println(data)

							s.Block.Scope.SetSymbolLocally("_", runtime.NewSymbol(result))
						}
					}
				}
			}
		}
	}
}
