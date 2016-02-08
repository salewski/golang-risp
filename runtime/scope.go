package runtime

type Scope struct {
	Symbols map[string]*Value
	parent  *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Symbols: make(map[string]*Value),
		parent:  parent,
	}
}
