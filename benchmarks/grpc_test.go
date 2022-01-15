package benchmarks

import (
	"context"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
)

func BenchmarkHTTP2GetWithWokers(b *testing.B) {
	wp := wpool.New(6)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx) //start workers

	requestQueue := wp.JobQueue()

	collector := wpool.NewCollector(wp.Results())

	go collector.Run(ctx) //start result collector

	job := &wpool.Job{
		ExecFn: func(ctx context.Context, args interface{}) ([]string, error) {
			return []string{"200_OK", "HTTP_1_1"}, nil
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		job.Details.ID = n
		requestQueue <- *job
	}
	close(requestQueue)
	<-wp.Done
	b.StopTimer()
	report := collector.GenerateReport()
	b.Log(report)
}
