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
		for key, _ := range s.Block.Scope.Symbols {
			if strings.HasPrefix(key, ident) {
				c = append(c, line[0:identEnd+1]+key)
			}
		}
	}

	return
}
