package compiler

type Pipeline interface {
	AddStage(stage Stage) Pipeline
	Execute(ctx *Context) error
}
