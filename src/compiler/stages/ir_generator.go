package stages

import (
	"github.com/kkumar-gcc/enumgen/src/compiler/ir"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type IRGenerator struct {
}

func NewIRGenerator() *IRGenerator {
	return &IRGenerator{}
}

func (r *IRGenerator) Name() string {
	return "IRGenerator"
}

func (r *IRGenerator) Process(ctx *compiler.Context) error {
	if ctx.AST == nil {
		return nil
	}

	transformer := ir.NewTransformer(ctx)
	ctx.IRModule = transformer.Transform()
	return nil
}
