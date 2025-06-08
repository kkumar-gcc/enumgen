package cmd

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

var rootCmd = &cli.Command{
	Name:  "edl",
	Usage: "A powerful tool for generating enum definitions",
	Description: `edl is a command-line tool designed to simplify the process of generating enum definitions in various programming languages.
It supports multiple languages and provides a flexible way to define enums using a simple syntax.`,
	Commands: []*cli.Command{
		generateCmd,
		langListCmd,
		langOptionsCmd,
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.Run(ctx, os.Args)
}
