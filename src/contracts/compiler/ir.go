package compiler

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type IRModule interface {
	Name() string
	Enums() []IREnumDefinition
	Source() string
	SetSource(source string) IRModule
	SetEnums(enums []IREnumDefinition) IRModule
	SetName(name string) IRModule
}

type IRValue interface {
	Position() token.Position
	String() string
}

type IRKeyValue interface {
	IRValue
	Key() IRValue
	Value() IRValue
}

type IRLiteral interface {
	IRValue
	Value() string
	Kind() token.Token
	TypeInfo() Type
}

type IREnumDefinition interface {
	Name() string
	Doc() string
	Members() []IREnumMember
	ValueType() Type
	KeyType() Type
	Position() token.Position
	OriginalNode() *ast.EnumDefinition
	FindMember(name string) IREnumMember
}

type IREnumMember interface {
	Name() string
	Doc() string
	Value() IRValue
	Position() token.Position
	OriginalNode() *ast.MemberDefinition
}
