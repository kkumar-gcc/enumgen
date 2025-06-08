package codegen

import (
	"fmt"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type CodeGenerationStage struct {
}

func NewCodeGenerationStage() *CodeGenerationStage {
	return &CodeGenerationStage{}
}

func (r *CodeGenerationStage) Name() string {
	return "CodeGeneration"
}

func (r *CodeGenerationStage) Process(ctx *compiler.Context) error {
	irModule := ctx.IRModule
	if irModule == nil {
		return fmt.Errorf("IR module generation failed")
	}

	if DefaultRegistry == nil {
		return fmt.Errorf("code generation registry is not initialized")
	}

	generator, err := DefaultRegistry.Get(ctx.TargetLang)
	if err != nil {
		fmt.Println("Error getting generator:", err)
		return err
	}

	files, err := generator.Generate(irModule, ctx.GenerationConfig)
	if err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	ctx.OutputFiles = files

	return nil
}
