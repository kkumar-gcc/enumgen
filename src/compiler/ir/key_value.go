package ir

import (
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type KeyValue struct {
	key   compiler.IRValue
	value compiler.IRValue
	pos   token.Position
}

var _ compiler.IRKeyValue = (*KeyValue)(nil)

func NewKeyValue(key compiler.IRValue, value compiler.IRValue, pos token.Position) *KeyValue {
	return &KeyValue{
		key:   key,
		value: value,
		pos:   pos,
	}
}

func (r *KeyValue) Position() token.Position {
	return r.pos
}

func (r *KeyValue) String() string {
	return r.key.String() + ": " + r.value.String()
}

func (r *KeyValue) Key() compiler.IRValue {
	return r.key
}

func (r *KeyValue) Value() compiler.IRValue {
	return r.value
}
