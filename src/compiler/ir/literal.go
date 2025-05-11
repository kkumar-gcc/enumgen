package ir

import (
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type Literal struct {
	kind     token.Token
	value    string
	pos      token.Position
	typeInfo compiler.Type
}

var _ compiler.IRValue = (*Literal)(nil)

func NewLiteral(kind token.Token, value string, pos token.Position, typeInfo compiler.Type) *Literal {
	return &Literal{
		kind:     kind,
		value:    value,
		pos:      pos,
		typeInfo: typeInfo,
	}
}

func (r *Literal) Position() token.Position {
	return r.pos
}

func (r *Literal) String() string {
	return r.value
}

func (r *Literal) Kind() token.Token {
	return r.kind
}

func (r *Literal) Value() string {
	return r.value
}

func (r *Literal) TypeInfo() compiler.Type {
	return r.typeInfo
}
