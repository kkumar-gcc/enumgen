package compiler

import (
	"fmt"
	"os"

	"github.com/kkumar-gcc/enumgen/src/codegen"
	"github.com/kkumar-gcc/enumgen/src/compiler/rules"
	"github.com/kkumar-gcc/enumgen/src/compiler/stages"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
)

var compilationRules = []compiler.Rule{
	rules.NewTypeCompatibilityRule(),
}

func CompileFile(filePath string, outputDir string, targetLang string, strict bool) (*compiler.Context, error) {
	options := make(map[string]any)

	ctx := &compiler.Context{
		SourcePath:       filePath,
		OutputDir:        outputDir,
		TargetLang:       targetLang,
		Errors:           make(errors.ErrorList, 0),
		GenerationConfig: options,
		Strict:           strict,
	}

	source, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	ctx.SourceCode = source

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	pipeline := NewPipeline()
	pipeline.AddStage(stages.NewParseStage()).
		AddStage(stages.NewSymbolCollector()).
		AddStage(stages.NewTypeResolver()).
		AddStage(stages.NewValidator(compilationRules)).
		AddStage(stages.NewIRGenerator()).
		AddStage(codegen.NewCodeGenerationStage())

	if err := pipeline.Execute(ctx); err != nil {
		return ctx, err
	}

	return ctx, nil
}
