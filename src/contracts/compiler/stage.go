package compiler

type Stage interface {
	Process(ctx *Context) error
	Name() string
}
