package types

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type CharFormatter struct{}

func (r *CharFormatter) GoTypeName() string { return "rune" }
func (r *CharFormatter) ZeroValue() any     { return rune(0) }
func (r *CharFormatter) FormatMemberValue(irValue compiler.IRValue, memberName string, index int) (any, error) {
	if irValue == nil {
		return rune(index), nil
	}

	literal, ok := irValue.(compiler.IRLiteral)
	if !ok {
		return nil, fmt.Errorf("internal error: expected IRLiteral for char handler, got %T", irValue)
	}

	if literal.Kind() != token.CHAR {
		return nil, fmt.Errorf("type error for member '%s': char enum expects a CHAR literal, got %v", memberName, literal.Kind())
	}

	unquoted, err := strconv.Unquote(literal.Value())
	if err != nil {
		return nil, fmt.Errorf("syntax error for member '%s': invalid character literal %s", memberName, literal.Value())
	}

	if utf8.RuneCountInString(unquoted) != 1 {
		return nil, fmt.Errorf("type error for member '%s': character literal %s must contain exactly one character", memberName, literal.Value())
	}

	charRune, _ := utf8.DecodeRuneInString(unquoted)

	return charRune, nil
}
