package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

var intTypeToBitSize = map[string]int{
	"int": 64, "int8": 8, "int16": 16, "int32": 32, "int64": 64,
	"uint": 64, "uint8": 8, "uint16": 16, "uint32": 32, "uint64": 64,
}

type IntFormatter struct {
	ConcreteGoType string
}

func (r *IntFormatter) GoTypeName() string { return r.ConcreteGoType }

func (r *IntFormatter) ZeroValue() any { return 0 }

func (r *IntFormatter) FormatMemberValue(irValue compiler.IRValue, memberName string, index int) (any, error) {
	numericStr, err := r.getNumericString(irValue)
	if err != nil {
		return nil, fmt.Errorf("for member '%s': %w", memberName, err)
	}

	if irValue == nil {
		numericStr = strconv.Itoa(index)
	}

	bitSize, ok := intTypeToBitSize[r.ConcreteGoType]
	if !ok {
		return nil, fmt.Errorf("internal error: unrecognized integer type '%s'", r.ConcreteGoType)
	}

	isUnsigned := strings.HasPrefix(r.ConcreteGoType, "u")
	if isUnsigned {
		val, err := strconv.ParseUint(numericStr, 0, bitSize)
		if err != nil {
			return nil, fmt.Errorf("invalid literal for member '%s': cannot parse '%s' as %s", memberName, numericStr, r.ConcreteGoType)
		}

		return val, nil
	}

	val, err := strconv.ParseInt(numericStr, 0, bitSize)
	if err != nil {
		return nil, fmt.Errorf("invalid literal for member '%s': cannot parse '%s' as %s", memberName, numericStr, r.ConcreteGoType)
	}

	return val, nil
}

func (r *IntFormatter) getNumericString(irValue compiler.IRValue) (string, error) {
	if irValue == nil {
		return "", nil
	}

	switch v := irValue.(type) {
	case compiler.IRLiteral:
		if v.Kind() != token.INT {
			return "", fmt.Errorf("type error: expected an INT literal, got %v", v.Kind())
		}
		return v.Value(), nil
	default:
		return "", fmt.Errorf("internal error: expected IRLiteral or IRUnary, got %T", irValue)
	}
}
