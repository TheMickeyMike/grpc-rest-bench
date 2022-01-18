package wpool

import (
	"context"
	"fmt"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

func New(wcount int) *WorkerPool {
	return &WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount*2),
		Done:         make(chan struct{}),
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results, i)
	}
	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp *WorkerPool) JobQueue() chan<- Job {
	return wp.jobs
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result, id int) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			select {
			case results <- job.execute(ctx): // Put results in the channel unless it is full
			default:
				fmt.Println("Channel full. Discarding value")
			}
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
