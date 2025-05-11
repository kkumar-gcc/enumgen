package compiler

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type SymbolTable interface {
	Define(symbol *Symbol) error
	Lookup(name string) *Symbol
	LookupEnum(name string) *Symbol
	CurrentScope() Scope
	SetCurrentScope(scope Scope)
	GlobalScope() Scope
	SetGlobalScope(scope Scope)
	EnterScope() Scope
	ExitScope() Scope
}

type Symbol struct {
	Name      string
	Kind      SymbolKind
	Node      ast.Node
	Type      Type
	Pos       token.Position
	Scope     Scope
	Docstring string
}

type SymbolKind int

const (
	SymbolEnum SymbolKind = iota
	SymbolEnumMember
	SymbolType
)

func (r SymbolKind) String() string {
	switch r {
	case SymbolEnum:
		return "enum"
	case SymbolEnumMember:
		return "enum member"
	case SymbolType:
		return "type"
	default:
		return "unknown"
	}
}
