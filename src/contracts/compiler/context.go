package compiler

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/errors"
)

type Context struct {
	SourcePath string
	SourceCode []byte

	AST *ast.File

	Symbols     SymbolTable
	Types       TypeRegistry
	Validations ValidationResult

	IRModule IRModule

	TargetLang       string
	OutputDir        string
	GenerationConfig map[string]string

	OutputFiles []*OutputFile

	Errors errors.ErrorList

	Strict bool
}

type OutputFile struct {
	Name string
	Path string
	Body []byte
}
