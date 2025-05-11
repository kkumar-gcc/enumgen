package compiler

import "github.com/kkumar-gcc/enumgen/src/ast"

type Type interface {
	Name() string
	String() string
	IsAssignableFrom(other Type) bool
	Kind() TypeKind
	Node() ast.Node
	EnumSymbol() *Symbol
	SetEnumSymbol(symbol *Symbol) Type
	ElementType() Type
	SetElementType(elementType Type) Type
	KeyType() Type
	SetKeyType(keyType Type) Type
	ValueType() Type
	SetValueType(valueType Type) Type
}

type TypeRegistry interface {
	RegisterType(t Type) error
	LookupType(name string) Type
	IsPrimitive(name string) bool
}

type TypeKind int

const (
	TypePrimitive TypeKind = iota
	TypeEnum
)
