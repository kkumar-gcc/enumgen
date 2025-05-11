package compiler

type Scope interface {
	Parent() Scope
	Children() []Scope
	SetChildren(children []Scope)
	Name() string
	Define(symbol *Symbol) error
	Lookup(name string) *Symbol
	LookupQualified(name string) *Symbol
	LookupLocal(name string) *Symbol
	Symbols() map[string]*Symbol
}
