package ir

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type EnumDefinition struct {
	name         string
	doc          string
	members      []compiler.IREnumMember
	valueType    compiler.Type
	keyType      compiler.Type
	position     token.Position
	originalNode *ast.EnumDefinition
}

func NewEnumDefinition(name string, doc string, members []compiler.IREnumMember, valueType compiler.Type, keyType compiler.Type, position token.Position, originalNode *ast.EnumDefinition) *EnumDefinition {
	return &EnumDefinition{
		name:         name,
		doc:          doc,
		members:      members,
		valueType:    valueType,
		keyType:      keyType,
		position:     position,
		originalNode: originalNode,
	}
}

func (r *EnumDefinition) Name() string {
	return r.name
}

func (r *EnumDefinition) Doc() string {
	return r.doc
}

func (r *EnumDefinition) Members() []compiler.IREnumMember {
	return r.members
}

func (r *EnumDefinition) ValueType() compiler.Type {
	return r.valueType
}

func (r *EnumDefinition) KeyType() compiler.Type {
	return r.keyType
}

func (r *EnumDefinition) Position() token.Position {
	return r.position
}

func (r *EnumDefinition) OriginalNode() *ast.EnumDefinition {
	return r.originalNode
}

func (r *EnumDefinition) String() string {
	return r.name + ": " + r.doc
}

func (r *EnumDefinition) FindMember(name string) compiler.IREnumMember {
	for _, member := range r.members {
		if member.Name() == name {
			return member
		}
	}
	return nil
}
