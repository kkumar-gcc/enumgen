package ir

import "github.com/kkumar-gcc/enumgen/src/contracts/compiler"

type Module struct {
	name   string
	enums  []compiler.IREnumDefinition
	source string
}

var _ compiler.IRModule = (*Module)(nil)

func NewModule(name string, source string) *Module {
	return &Module{
		name:   name,
		enums:  []compiler.IREnumDefinition{},
		source: source,
	}
}

func (r *Module) Name() string {
	return r.name
}

func (r *Module) Enums() []compiler.IREnumDefinition {
	return r.enums
}

func (r *Module) Source() string {
	return r.source
}

func (r *Module) SetSource(source string) compiler.IRModule {
	r.source = source
	return r
}

func (r *Module) SetEnums(enums []compiler.IREnumDefinition) compiler.IRModule {
	r.enums = enums
	return r
}

func (r *Module) SetName(name string) compiler.IRModule {
	r.name = name
	return r
}
