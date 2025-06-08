package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/kkumar-gcc/enumgen/src/codegen"
)

var langOptionsCmd = &cli.Command{
	Name:  "lang-options",
	Usage: "List available options for a specific language",
	Description: `The lang-options command displays all available options for a specific programming language used in enum generation.
It provides detailed information about each option, including its default value and a brief description.`,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name: "lang",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		lang := cmd.StringArg("lang")
		if lang == "" {
			return cli.Exit("Error: No language specified. Please provide a programming language to list options.", 1)
		}

		generator, err := codegen.DefaultRegistry.Get(lang)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Error: Language '%s' not found. Please check the available languages.", lang), 1)
		}

		options := generator.DefaultOptions()
		if len(options) == 0 {
			return cli.Exit(fmt.Sprintf("No options available for language '%s'.", lang), 0)
		}

		fmt.Printf("Available options for language '%s':\n", lang)
		fmt.Println(generator.OptionHelp())

		return nil
	},
}
