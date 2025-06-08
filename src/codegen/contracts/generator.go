package contracts

import (
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Generator interface {
	Name() string
	Language() string
	FileExtension() string
	Generate(module compiler.IRModule, options map[string]string) ([]*compiler.OutputFile, error)
	GetDefaultOptions() map[string]string
	OptionHelp() string
}
