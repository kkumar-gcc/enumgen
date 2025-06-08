package compiler

import (
	"fmt"
	"os"

	"github.com/kkumar-gcc/enumgen/src/compiler/stages"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
)

func ParseFile(filePath string) (*compiler.Context, error) {
	ctx := &compiler.Context{
		SourcePath: filePath,
		Errors:     make(errors.ErrorList, 0),
	}

	source, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	ctx.SourceCode = source

	pipeline := NewPipeline()
	pipeline.AddStage(stages.NewParseStage())
	if err := pipeline.Execute(ctx); err != nil {
		return ctx, err
	}

	return ctx, nil
}
