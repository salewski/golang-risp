package lexer

import (
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	pos      int
	startPos int
	line     int
	col      int
	data     string
	source   Source
	Tokens   []*Token
}

func NewLexer(source Source) *Lexer {
	return &Lexer{
		pos:      0,
		startPos: 0,
		line:     1,
		col:      1,
		data:     source.Data(),
		source:   source,
	}
}

func (l *Lexer) addToken(typ TokenType) *Token {
	token := NewToken(typ, l.buffer(), l.newPos())

	l.Tokens = append(l.Tokens, token)

	l.resetBuffer()

	return token
}

func (l *Lexer) newPos() *TokenPos {
	return &TokenPos{
		Line:   l.line,
		Col:    l.col,
		Source: l.source,
	}
}

func (l *Lexer) current() rune {
	return l.peek(0)
}

func (l *Lexer) peek(amount int) rune {
	return rune(l.data[l.pos+amount])
}

func (l *Lexer) buffer() string {
	return string(l.data[l.startPos:l.pos])
}

func (l *Lexer) resetBuffer() {
	l.startPos = l.pos
}

func (l *Lexer) consume() {
	if l.current() == '\n' {
		l.col = 1
		l.line++
	} else {
		l.col++
	}

	l.pos++
}

func (l *Lexer) ignore(amount int) {
	for i := 0; i < amount; i++ {
		l.consume()
	}

	l.resetBuffer()
}

func (l *Lexer) next() {
	l.pos++
}

func (l *Lexer) isEOF() bool {
	return l.pos >= len(l.data)
}

func (l *Lexer) hasNext() bool {
	return l.pos < len(l.data)
}

func (l *Lexer) lexNumber() error {
	if l.current() == '+' || l.current() == '-' {
		l.consume()
	}

	for l.hasNext() && isNumber(l.current()) {
		l.consume()

		if !l.isEOF() && l.current() == '.' {
			if strings.Contains(l.buffer(), ".") {
				return NewSyntaxError(l.newPos(), "malformed number")
			}

			l.consume()
		}
	}

	l.addToken(Number)

	return nil
}

func (l *Lexer) lexString() error {
	l.ignore(1)

	for l.hasNext() && l.current() != '"' {
		l.consume()
	}

	if l.isEOF() || l.current() != '"' {
		return NewSyntaxError(l.newPos(), "unclosed string literal")
	}

	t := l.addToken(String)
	t.Data, _ = strconv.Unquote("\"" + t.Data + "\"")

	l.ignore(1)

	return nil
}

func (l *Lexer) lexIdentifierOrKeyword(keyword bool) {
	typ := Identifier

	if keyword {
		typ = Keyword

		l.ignore(1)
	}

	l.consume()

	for l.hasNext() && isIdentifierPart(l.current()) {
		l.consume()
	}

	l.addToken(typ)
}

func (l *Lexer) lexSeparator() {
	l.consume()
	l.addToken(Separator)
}

func (l *Lexer) Lex() error {
	for !l.isEOF() {
		switch {
		case isNumber(l.current()), (l.current() == '+' || l.current() == '-') && l.hasNext() && isNumber(l.peek(1)):
			err := l.lexNumber()

			if err != nil {
				return err
			}
		case (l.current() == '>' || l.current() == '<' || l.current() == '!') && l.hasNext() && l.peek(1) == '=':
			l.consume()
			l.consume()
			l.addToken(Identifier)
		case l.current() == '_', l.current() == '+', l.current() == '-', l.current() == '*', l.current() == '/', l.current() == '%', l.current() == '=', l.current() == '>', l.current() == '<':
			l.consume()
			l.addToken(Identifier)
		case l.current() == ':' && l.hasNext() && isIdentifierStart(l.peek(1)):
			l.lexIdentifierOrKeyword(true)
		case isIdentifierStart(l.current()):
			l.lexIdentifierOrKeyword(false)
		case l.current() == '"':
			err := l.lexString()

			if err != nil {
				return err
			}
		case l.current() == '(', l.current() == ')':
			l.lexSeparator()
		case l.current() < ' ', unicode.IsControl(l.current()), unicode.IsSpace(l.current()):
			l.ignore(1)
		default:
			return NewSyntaxError(l.newPos(), "unknown character '%c'", l.current())
		}
	}

	return nil
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r)
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || unicode.IsNumber(r) || r == '-' || r == ':'
}
