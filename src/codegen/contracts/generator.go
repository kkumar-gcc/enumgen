package contracts

import (
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Generator interface {
	Name() string
	Language() string
	DefaultOptions() map[string]string
	OptionHelp() string
	Generate(module compiler.IRModule, options map[string]string) ([]*compiler.OutputFile, error)
}
