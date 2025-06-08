package types

import (
	"fmt"
	"strconv"

	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type FloatHandler struct {
	ConcreteGoType string
}

func (r *FloatHandler) GoTypeName() string { return r.ConcreteGoType }
func (r *FloatHandler) ZeroValue() any {
	return 0.0
}

func (r *FloatHandler) FormatMemberValue(irValue compiler.IRValue, memberName string, index int) (any, error) {
	bitSize := 64
	if r.ConcreteGoType == "float32" {
		bitSize = 32
	}

	numericStr, kind, err := r.getNumericString(irValue)
	if err != nil {
		return nil, fmt.Errorf("for member '%s': %w", memberName, err)
	}

	if kind != token.INT && kind != token.FLOAT {
		return nil, fmt.Errorf("type error for member '%s': float enum expects INT or FLOAT literal, got %v", memberName, kind)
	}

	val, err := strconv.ParseFloat(numericStr, bitSize)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric literal for member '%s': cannot parse '%s' as %s", memberName, numericStr, r.ConcreteGoType)
	}

	return val, nil
}

func (r *FloatHandler) getNumericString(irValue compiler.IRValue) (valueStr string, kind token.Token, err error) {
	switch v := irValue.(type) {
	case compiler.IRLiteral:
		return v.Value(), v.Kind(), nil
	default:
		return "", token.ILLEGAL, fmt.Errorf("internal error: expected IRLiteral or IRUnary for float handler, got %T", irValue)
	}
}
