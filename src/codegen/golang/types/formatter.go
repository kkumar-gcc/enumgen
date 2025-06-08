package types

import "github.com/kkumar-gcc/enumgen/src/contracts/compiler"

type ValueFormatter interface {
	// GoTypeName returns the actual Go type name (e.g., "string", "int32").
	GoTypeName() string

	// ZeroValue returns the zero value for this Go type (e.g., `""`, 0).
	ZeroValue() any

	// FormatMemberValue processes an IRLiteral to its Go code representation.
	// The `index` is used for default/iota-like values when irValue is nil.
	FormatMemberValue(irValue compiler.IRValue, memberName string, index int) (any, error)
}
