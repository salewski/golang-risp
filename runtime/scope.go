package runtime

type Symtab map[string]*Symbol
type Mactab map[string]*Macro

type Scope struct {
	Symbols Symtab
	Macros  Mactab
	parent  *Scope
}

type Symbol struct {
	Value    *Value
	Exported bool
}

func NewSymbol(value *Value) *Symbol {
	return &Symbol{
		Value:    value,
		Exported: false,
	}
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Symbols: make(Symtab),
		Macros:  make(Mactab),
		parent:  parent,
	}
}

func (s *Scope) ApplySymbols(symbols Symtab) {
	for key, value := range symbols {
		s.Symbols[key] = value
	}
}

func (s *Scope) GetSymbol(key string) *Symbol {
	if s.Symbols[key] == nil && s.parent != nil {
		return s.parent.GetSymbol(key)
	}

	return s.Symbols[key]
}

func (s *Scope) SetSymbol(key string, value *Symbol) {
	s.Symbols[key] = value
}

func (s *Scope) RemoveSymbol(key string) {
	delete(s.Symbols, key)
}

func (s *Scope) HasSymbol(key string) bool {
	return s.GetSymbol(key) != nil
}

func (s *Scope) ApplyMacros(macros Mactab) {
	for key, value := range macros {
		s.Macros[key] = value
	}
}

func (s *Scope) GetMacro(key string) *Macro {
	if s.Macros[key] == nil && s.parent != nil {
		return s.parent.GetMacro(key)
	}

	return s.Macros[key]
}

func (s *Scope) HasMacro(key string) bool {
	return s.GetMacro(key) != nil
}
