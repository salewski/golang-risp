package repl

import (
	"fmt"
	"github.com/peterh/liner"
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"github.com/raoulvdberge/risp/util"
	"os"
	"path/filepath"
	"strconv"
)

var (
	history = filepath.Join(os.TempDir(), ".risp_repl")
)

type ReplSession struct {
	block  *runtime.Block
	tokens []*lexer.Token
	depth  int
}

func NewReplSession(block *runtime.Block) *ReplSession {
	return &ReplSession{
		block: block,
		depth: 0,
	}
}

func (s *ReplSession) Run() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(s.completer)

	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

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

		if f, err := os.Create(history); err == nil {
			line.WriteHistory(f)
			f.Close()
		}

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

					continue
				}

				for _, node := range p.Nodes {
					result, err := s.block.EvalNode(node)

					if err != nil {
						util.ReportError(err, true)

						continue
					}

					data := util.Yellow("===> " + result.String())

					if resultType(result) != "" {
						data += " " + util.Yellow("("+resultType(result)+")")
					}

					fmt.Println(data)

					s.block.Scope.SetSymbolLocally("_", runtime.NewSymbol(result))
				}
			}
		}
	}
}

func resultType(value *runtime.Value) string {
	switch value.Type {
	case runtime.NilValue:
		return ""
	case runtime.QuotedValue:
		return "quoted " + value.Quoted.Name()
	default:
		return value.Type.String()
	}
}
