package wpool

import "context"

type ExecutionFn func(ctx context.Context, args interface{}) ([]string, error)

type JobDetails struct {
	ID       int
	Metadata map[string]string
}

type Result struct {
	Value   []string
	Err     error
	Details JobDetails
}

type Job struct {
	Details JobDetails
	ExecFn  ExecutionFn
	Args    interface{}
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:     err,
			Details: j.Details,
		}
	}

	return Result{
		Value:   value,
		Details: j.Details,
	}
}
