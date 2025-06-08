package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/kkumar-gcc/enumgen/src/codegen"
)

var langListCmd = &cli.Command{
	Name:  "lang-list",
	Usage: "List all available languages for enum generation",
	Description: `The lang-list command displays all the programming languages supported for enum generation in the tool.
It provides a quick overview of the languages available, allowing users to choose the appropriate one for their enum definitions.`,
	Action: func(ctx context.Context, cmd *cli.Command) error {
		languages := codegen.DefaultRegistry.Languages()
		for _, lang := range languages {
			fmt.Println(lang)
		}
		return nil
	},
}
