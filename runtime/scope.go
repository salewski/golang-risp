package runtime

type Symtab map[string]*Value

type Scope struct {
	Symbols Symtab
	parent  *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Symbols: make(Symtab),
		parent:  parent,
	}
}

func (s *Scope) Apply(symbols Symtab) {
	for key, value := range symbols {
		s.Symbols[key] = value
	}
}
