package contracts

import (
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Generator interface {
	Name() string
	Language() string
	FileExtension() string
	Generate(module compiler.IRModule, outputDir string, options map[string]interface{}) ([]*compiler.OutputFile, error)
	GetDefaultOptions() map[string]interface{}
	OptionHelp() string
}
