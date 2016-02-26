package repl

import (
	"github.com/raoulvdberge/risp/lexer"
	"github.com/raoulvdberge/risp/util"
	"strings"
)

func (s *ReplSession) completer(line string) (c []string) {
	identEnd := -1
	ident := ""

	for i := len(line) - 1; i >= 0; i-- {
		part := rune(line[i])

		if lexer.IsIdentifierPart(part) {
			identEnd = i
			ident += string(part)
		} else if lexer.IsIdentifierStart(part) {
			identEnd = i
			ident += string(part)

			break
		} else {
			break
		}
	}

	ident = util.ReverseString(ident)

	if ident != "" && lexer.IsIdentifierStart(rune(ident[0])) && identEnd != -1 {
		prev := line[0:identEnd]

		for name, _ := range s.block.Scope.Symbols {
			if strings.HasPrefix(name, ident) {
				c = append(c, prev+name)
			}
		}

		for name, _ := range s.block.Scope.Macros {
			if strings.HasPrefix(name, ident) {
				c = append(c, prev+name)
			}
		}
	}

	return
}
