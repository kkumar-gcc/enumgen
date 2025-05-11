package compiler

import "github.com/kkumar-gcc/enumgen/src/contracts/compiler"

type Pipeline struct {
	stages []compiler.Stage
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		stages: make([]compiler.Stage, 0),
	}
}

func (r *Pipeline) AddStage(stage compiler.Stage) compiler.Pipeline {
	r.stages = append(r.stages, stage)
	return r
}

func (r *Pipeline) Execute(ctx *compiler.Context) error {
	for _, stage := range r.stages {
		if err := stage.Process(ctx); err != nil {
			return err
		}

		if !ctx.Errors.IsEmpty() && ctx.Errors.HasFatal() {
			return ctx.Errors
		}
	}

	return nil
}
