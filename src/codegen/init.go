package codegen

import (
	"sync"

	"github.com/kkumar-gcc/enumgen/src/codegen/golang"
)

var (
	DefaultRegistry *Registry
	once            sync.Once
)

func Init() {
	once.Do(func() {
		DefaultRegistry = NewRegistry()

		goGenerator, err := golang.New()
		if err != nil {
			panic("failed to initialize Go generator: " + err.Error())
		}
		DefaultRegistry.Register(goGenerator)
	})
}
