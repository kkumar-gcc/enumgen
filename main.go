package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kkumar-gcc/enumgen/cmd"
	"github.com/kkumar-gcc/enumgen/src/codegen"
)

func main() {
	codegen.Init()

	if err := cmd.Execute(context.Background()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error executing command: %v\n\n", err)
		os.Exit(1)
	}
}
