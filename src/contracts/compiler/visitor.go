package compiler

import "github.com/kkumar-gcc/enumgen/src/ast"

type Visitor interface {
	VisitFile(node *ast.File) any
	VisitEnum(node *ast.EnumDefinition) any
	VisitTypeSpec(node *ast.TypeSpec) any
	VisitTypeRef(node *ast.TypeRef) any
	VisitMember(node *ast.MemberDefinition) any
	VisitBasicLit(node *ast.BasicLit) any
	VisitKeyValueExpr(node *ast.KeyValueExpr) any
	VisitIdent(node *ast.Ident) any
}
