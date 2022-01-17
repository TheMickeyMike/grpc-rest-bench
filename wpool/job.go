package wpool

import "context"

type ExecutionFn func(ctx context.Context) (string, error)

type Result struct {
	Value string
	Err   error
}

type Job struct {
	ExecFn ExecutionFn
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx)
	if err != nil {
		return Result{
			Err: err,
		}
	}

	return Result{
		Value: value,
	}
}
