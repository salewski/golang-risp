package runtime

type Symtab map[string]*Value

type Scope struct {
	symbols Symtab
	parent  *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		symbols: make(Symtab),
		parent:  parent,
	}
}

func (s *Scope) Apply(symbols Symtab) {
	for key, value := range symbols {
		s.symbols[key] = value
	}
}

func (s *Scope) Get(key string) *Value {
	if s.symbols[key] == nil && s.parent != nil {
		return s.parent.Get(key)
	}

	return s.symbols[key]
}

func (s *Scope) Has(key string) bool {
	return s.Get(key) != nil
}
