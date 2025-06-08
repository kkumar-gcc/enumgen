package types

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/pkg/strconvx"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type StringFormatter struct{}

func (r *StringFormatter) GoTypeName() string { return "string" }
func (r *StringFormatter) ZeroValue() any     { return "" }

func (r *StringFormatter) FormatMemberValue(irValue compiler.IRValue, memberName string, index int) (any, error) {
	if irValue == nil {
		return memberName, nil
	}

	literal, ok := irValue.(compiler.IRLiteral)
	if !ok {
		return nil, fmt.Errorf("internal error: expected IRLiteral for string handler, got %T", irValue)
	}

	if literal.Kind() != token.STRING {
		return nil, fmt.Errorf("type error for member '%s': string enum expects a STRING literal, got %v", memberName, literal.Kind())
	}

	return strconvx.Unquote(literal.Value()), nil
}
