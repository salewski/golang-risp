package runtime

type Symtab map[string]*Value
type Mactab map[string]*Macro

type Scope struct {
	symbols Symtab
	macros  Mactab
	parent  *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		symbols: make(Symtab),
		macros:  make(Mactab),
		parent:  parent,
	}
}

func (s *Scope) IsOuterMost() bool {
	return s.parent == nil
}

func (s *Scope) ApplySymbols(symbols Symtab) {
	for key, value := range symbols {
		s.symbols[key] = value
	}
}

func (s *Scope) GetSymbol(key string) *Value {
	if s.symbols[key] == nil && s.parent != nil {
		return s.parent.GetSymbol(key)
	}

	return s.symbols[key]
}

func (s *Scope) SetSymbol(key string, value *Value) {
	s.symbols[key] = value
}

func (s *Scope) RemoveSymbol(key string) {
	delete(s.symbols, key)
}

func (s *Scope) HasSymbol(key string) bool {
	return s.GetSymbol(key) != nil
}

func (s *Scope) ApplyMacros(macros Mactab) {
	for key, value := range macros {
		s.macros[key] = value
	}
}

func (s *Scope) GetMacro(key string) *Macro {
	if s.macros[key] == nil && s.parent != nil {
		return s.parent.GetMacro(key)
	}

	return s.macros[key]
}

func (s *Scope) HasMacro(key string) bool {
	return s.GetMacro(key) != nil
}
