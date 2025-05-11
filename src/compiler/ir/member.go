package ir

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type EnumMember struct {
	name         string
	doc          string
	value        compiler.IRValue
	position     token.Position
	originalNode *ast.MemberDefinition
}

func NewEnumMember(name string, doc string, value compiler.IRValue, position token.Position, originalNode *ast.MemberDefinition) *EnumMember {
	return &EnumMember{
		name:         name,
		doc:          doc,
		value:        value,
		position:     position,
		originalNode: originalNode,
	}
}

func (r *EnumMember) Name() string {
	return r.name
}

func (r *EnumMember) Doc() string {
	return r.doc
}

func (r *EnumMember) Value() compiler.IRValue {
	return r.value
}

func (r *EnumMember) Position() token.Position {
	return r.position
}

func (r *EnumMember) OriginalNode() *ast.MemberDefinition {
	return r.originalNode
}

func (r *EnumMember) String() string {
	return r.name + ": " + r.value.String()
}
