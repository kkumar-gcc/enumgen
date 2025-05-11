package symbols

import "github.com/kkumar-gcc/enumgen/src/contracts/compiler"

type Table struct {
	globalScope  compiler.Scope
	currentScope compiler.Scope

	enumsByName map[string]*compiler.Symbol
	types       map[string]compiler.Type
}

func NewTable() *Table {
	global := NewScope(nil)
	return &Table{
		globalScope:  global,
		currentScope: global,
		enumsByName:  make(map[string]*compiler.Symbol),
		types:        make(map[string]compiler.Type),
	}
}

func (r *Table) CurrentScope() compiler.Scope {
	return r.currentScope
}

func (r *Table) SetCurrentScope(scope compiler.Scope) {
	r.currentScope = scope
}

func (r *Table) GlobalScope() compiler.Scope {
	return r.globalScope
}

func (r *Table) SetGlobalScope(scope compiler.Scope) {
	r.globalScope = scope
}

func (r *Table) EnterScope() compiler.Scope {
	r.currentScope = NewScope(r.currentScope)
	return r.currentScope
}

func (r *Table) ExitScope() compiler.Scope {
	if r.currentScope.Parent() != nil {
		r.currentScope = r.currentScope.Parent()
	}

	return r.currentScope
}

func (r *Table) Define(symbol *compiler.Symbol) error {
	if err := r.currentScope.Define(symbol); err != nil {
		return err
	}
	if symbol.Kind == compiler.SymbolEnum {
		r.enumsByName[symbol.Name] = symbol
	}
	return nil
}

func (r *Table) Lookup(name string) *compiler.Symbol {
	return r.currentScope.Lookup(name)
}

func (r *Table) LookupEnum(name string) *compiler.Symbol {
	if symbol, ok := r.enumsByName[name]; ok {
		return symbol
	}
	return nil
}
