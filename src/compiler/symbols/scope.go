package symbols

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Scope struct {
	parent   compiler.Scope
	children []compiler.Scope
	name     string
	symbols  map[string]*compiler.Symbol
}

var _ compiler.Scope = (*Scope)(nil)

func NewScope(parent compiler.Scope) *Scope {
	s := &Scope{
		parent:   parent,
		children: []compiler.Scope{},
		name:     "",
		symbols:  make(map[string]*compiler.Symbol),
	}

	if parent != nil {
		parent.SetChildren([]compiler.Scope{s})
	}

	return s
}

func (r *Scope) Parent() compiler.Scope {
	return r.parent
}

func (r *Scope) Children() []compiler.Scope {
	return r.children
}

func (r *Scope) SetChildren(children []compiler.Scope) {
	r.children = append(r.children, children...)
}

func (r *Scope) Name() string {
	return r.name
}

func (r *Scope) Define(symbol *compiler.Symbol) error {
	name := symbol.Name
	if existing, ok := r.symbols[name]; ok {
		return fmt.Errorf("symbol %q already defined at %v", name, existing.Pos)
	}
	r.symbols[name] = symbol
	symbol.Scope = r
	return nil
}

func (r *Scope) Lookup(name string) *compiler.Symbol {
	if symbol, ok := r.symbols[name]; ok {
		return symbol
	}

	if r.parent != nil {
		return r.parent.Lookup(name)
	}

	return nil
}

func (r *Scope) LookupQualified(name string) *compiler.Symbol {
	return r.Lookup(name)
}

func (r *Scope) LookupLocal(name string) *compiler.Symbol {
	if symbol, ok := r.symbols[name]; ok {
		return symbol
	}

	return nil
}

func (r *Scope) Symbols() map[string]*compiler.Symbol {
	return r.symbols
}
