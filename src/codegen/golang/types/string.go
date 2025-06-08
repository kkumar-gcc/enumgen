package types

import (
	"fmt"
	"strconv"

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

	unquoted, err := strconv.Unquote(literal.Value())
	if err != nil {
		return nil, fmt.Errorf("syntax error for member '%s': invalid string literal %s", memberName, literal.Value())
	}

	return unquoted, nil
}
