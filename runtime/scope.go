package runtime

type Symtab map[string]*Symbol
type Mactab map[string]*Macro

type Scope struct {
	symbols Symtab
	macros  Mactab
	parent  *Scope
}

type Symbol struct {
	Value *Value
}

func NewSymbol(value *Value) *Symbol {
	return &Symbol{Value: value}
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		symbols: make(Symtab),
		macros:  make(Mactab),
		parent:  parent,
	}
}

func (s *Scope) ApplySymbols(symbols Symtab) {
	for key, value := range symbols {
		s.symbols[key] = value
	}
}

func (s *Scope) GetSymbol(key string) *Symbol {
	if s.symbols[key] == nil && s.parent != nil {
		return s.parent.GetSymbol(key)
	}

	return s.symbols[key]
}

func (s *Scope) SetSymbol(key string, value *Symbol) {
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
