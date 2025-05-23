package types

import (
	"fmt"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Registry struct {
	// Built-in primitive types
	primitives map[string]compiler.Type

	// All types by name
	types map[string]compiler.Type
}

func NewRegistry() *Registry {
	return &Registry{
		primitives: primitives,
		types:      primitives,
	}
}

func (r *Registry) RegisterType(t compiler.Type) error {
	name := t.Name()
	if existing, ok := r.types[name]; ok {
		return fmt.Errorf("type %q already registered as %v", name, existing.Kind())
	}

	r.types[name] = t
	return nil
}

func (r *Registry) LookupType(name string) compiler.Type {
	if t, ok := r.types[name]; ok {
		return t
	}
	return nil
}

func (r *Registry) IsPrimitive(name string) bool {
	_, ok := r.primitives[name]
	return ok
}

var primitives = map[string]compiler.Type{
	"int":     NewType(compiler.TypePrimitive, "int", nil),
	"string":  NewType(compiler.TypePrimitive, "string", nil),
	"bool":    NewType(compiler.TypePrimitive, "bool", nil),
	"float":   NewType(compiler.TypePrimitive, "float", nil),
	"int32":   NewType(compiler.TypePrimitive, "int32", nil),
	"int64":   NewType(compiler.TypePrimitive, "int64", nil),
	"uint":    NewType(compiler.TypePrimitive, "uint", nil),
	"uint32":  NewType(compiler.TypePrimitive, "uint32", nil),
	"uint64":  NewType(compiler.TypePrimitive, "uint64", nil),
	"float32": NewType(compiler.TypePrimitive, "float32", nil),
	"float64": NewType(compiler.TypePrimitive, "float64", nil),
}

func IsPrimitiveType(typeName string) bool {
	_, ok := primitives[typeName]
	return ok
}
