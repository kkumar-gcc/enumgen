package cmd

import (
	"context"
	"github.com/urfave/cli/v3"

	"github.com/kkumar-gcc/enumgen/src/compiler"
)

var generateCmd = &cli.Command{
	Name:  "generate",
	Usage: "Generate enum definitions from source files",
	Description: `The generate command processes source files to produce enum definitions in the specified output format.
It reads the source files, parses them, and generates the corresponding enum definitions based on the provided specifications.`,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name: "file",
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output directory for generated files",
		},
		&cli.StringFlag{
			Name:    "lang",
			Aliases: []string{"l"},
			Usage:   "Target programming language for enum definitions",
			Value:   "go",
		},
		&cli.BoolFlag{
			Name:    "strict",
			Aliases: []string{"s"},
			Usage:   "Enable strict mode for validation",
			Value:   false,
		},
		&cli.StringMapFlag{
			Name:    "options",
			Aliases: []string{"O"},
			Usage:   "Additional generation options in key=value format",
			Value:   make(map[string]string),
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fileName := cmd.StringArg("file")
		if fileName == "" {
			return cli.Exit("Error: No file specified. Please provide a source file to generate enum definitions.", 1)
		}

		outputDir := cmd.String("output")
		targetLang := cmd.String("lang")
		strict := cmd.Bool("strict")
		generationOptions := cmd.StringMap("options")
		compilerCtx, err := compiler.CompileFile(fileName, outputDir, targetLang, strict, generationOptions)
		if err != nil {
			return cli.Exit("Error: Failed to compile file: "+err.Error(), 1)
		}
		if len(compilerCtx.Errors) > 0 {
			return cli.Exit("Error: Compilation errors occurred:\n"+compilerCtx.Errors.Summary(), 1)
		}
		return nil
	},
}
