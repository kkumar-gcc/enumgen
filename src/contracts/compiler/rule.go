package compiler

import "github.com/kkumar-gcc/enumgen/src/ast"

type Rule interface {
	Name() string
	Check(ctx *Context, node ast.Node) []Issue
}
