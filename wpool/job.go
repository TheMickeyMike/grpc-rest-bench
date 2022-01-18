package wpool

import (
	"context"
)

type ExecutionFn func(ctx context.Context) (string, int64, error)

type Result struct {
	Value   string
	Retries int64
	Err     error
}

type Job struct {
	ExecFn ExecutionFn
}

func (j *Job) execute(ctx context.Context) Result {
	value, retries, err := j.ExecFn(ctx)
	if err != nil {
		return Result{
			Err: err,
		}
	}

	return Result{
		Value:   value,
		Retries: retries,
	}
}
