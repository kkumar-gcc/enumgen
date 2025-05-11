package types

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Type struct {
	kind compiler.TypeKind
	name string
	node ast.Node

	// For enums
	enumSymbol *compiler.Symbol
	// For container types
	elementType compiler.Type
	keyType     compiler.Type
	valueType   compiler.Type
}

func NewType(kind compiler.TypeKind, name string, node ast.Node) *Type {
	return &Type{
		kind: kind,
		name: name,
		node: node,
	}
}

func (t *Type) Kind() compiler.TypeKind {
	return t.kind
}

func (t *Type) Name() string {
	return t.name
}

func (t *Type) Node() ast.Node {
	return t.node
}

func (t *Type) EnumSymbol() *compiler.Symbol {
	return t.enumSymbol
}

func (t *Type) SetEnumSymbol(symbol *compiler.Symbol) compiler.Type {
	t.enumSymbol = symbol
	return t
}

func (t *Type) ElementType() compiler.Type {
	return t.elementType
}

func (t *Type) SetElementType(elementType compiler.Type) compiler.Type {
	t.elementType = elementType
	return t
}

func (t *Type) KeyType() compiler.Type {
	return t.keyType
}

func (t *Type) SetKeyType(keyType compiler.Type) compiler.Type {
	t.keyType = keyType
	return t
}

func (t *Type) ValueType() compiler.Type {
	return t.valueType
}

func (t *Type) SetValueType(valueType compiler.Type) compiler.Type {
	t.valueType = valueType
	return t
}

func (t *Type) String() string {
	if t.kind == compiler.TypeEnum {
		return t.enumSymbol.Name
	}
	return t.name
}

func (t *Type) IsAssignableFrom(other compiler.Type) bool {
	if t.kind == other.Kind() {
		return true
	}
	if t.kind == compiler.TypeEnum && other.Kind() == compiler.TypePrimitive {
		return true
	}
	if t.kind == compiler.TypePrimitive && other.Kind() == compiler.TypeEnum {
		return true
	}
	return false
}
