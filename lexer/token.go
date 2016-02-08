package lexer

type TokenType int

const (
	Number TokenType = iota
	String
	Identifier
	Keyword
	Separator
)

type Token struct {
	Type TokenType `json:"type"`
	Data string    `json:"data"`
	Pos  *TokenPos `json:"pos"`
}

type TokenPos struct {
	Line   int    `json:"line"`
	Col    int    `json:"col"`
	Source Source `json:"-"`
}

func NewToken(typ TokenType, data string, pos *TokenPos) *Token {
	return &Token{Type: typ, Data: data, Pos: pos}
}

func (t *Token) IsType(typ TokenType) bool {
	return t.Type == typ
}

func (t *Token) IsTypeAndData(typ TokenType, data string) bool {
	return t.IsType(typ) && t.Data == data
}

func (t *Token) DepthModifier() int {
	if t.IsTypeAndData(Separator, "(") {
		return 1
	} else if t.IsTypeAndData(Separator, ")") {
		return -1
	}
	return 0
}
