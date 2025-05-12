package codegen

import (
	"sync"
)

var (
	DefaultRegistry *Registry
	once            sync.Once
)

func Init() {
	once.Do(func() {
		DefaultRegistry = NewRegistry()

		//DefaultRegistry.Register(golang.NewGenerator())
		//DefaultRegistry.Register(typescript.NewGenerator())
	})
}
