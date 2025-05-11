package stages

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
	"github.com/kkumar-gcc/enumgen/src/lexer"
	"github.com/kkumar-gcc/enumgen/src/parser"
)

type ParseStage struct {
}

func NewParseStage() *ParseStage {
	return &ParseStage{}
}

func (r *ParseStage) Name() string {
	return "Parse"
}

func (r *ParseStage) Process(ctx *compiler.Context) error {
	lex := lexer.New(ctx.SourcePath, ctx.SourceCode, 0)

	p := parser.New(lex)
	file := p.Parse()

	if errs := p.Errors(); len(errs) > 0 {
		for _, err := range errs {
			ctx.Errors.Add(&errors.CompilationError{
				Pos:      err.Pos,
				Msg:      err.Msg,
				Severity: errors.SeverityError,
				Stage:    r.Name(),
				Filename: ctx.SourcePath,
			})
		}
		return fmt.Errorf("parse errors: %v", errs)
	}

	ctx.AST = file
	return nil
}
