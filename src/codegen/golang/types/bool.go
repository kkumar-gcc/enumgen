package types

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type BoolHandler struct{}

func (r *BoolHandler) GoTypeName() string { return "bool" }
func (r *BoolHandler) ZeroValue() any     { return false }

func (r *BoolHandler) FormatMemberValue(irValue compiler.IRValue, memberName string, index int) (any, error) {
	if irValue == nil {
		return false, nil
	}

	literal, ok := irValue.(compiler.IRLiteral)
	if !ok {
		return nil, fmt.Errorf("internal error: expected IRLiteral for bool handler, got %T", irValue)
	}

	switch kind := literal.Kind(); kind {
	case token.TRUE:
		return true, nil
	case token.FALSE:
		return false, nil
	default:
		pos := literal.Position()
		return nil, fmt.Errorf(
			"type error at %s member '%s' expects a boolean (true or false), but got %v",
			pos.String(),
			memberName,
			kind,
		)
	}
}
