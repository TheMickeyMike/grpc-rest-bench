package benchmarks

import (
	"context"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
)

func BenchmarkHTTP2GetWithWokers(b *testing.B) {
	client := NewHTTPClient(HTTP)

	wp := wpool.New(1)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx) //start workers

	requestQueue := wp.JobQueue()

	collector := wpool.NewCollector(wp.Results())

	go collector.Run(ctx) //start result collector

	job := wpool.Job{
		ExecFn: func(ctx context.Context) (string, error) {
			var (
				err      error
				response ResponseDetails
				target   pb.UserAccount
				retries  int = 3
			)
			for retries > 0 {
				response, err = client.MakeRequest(ctx, "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c", &target)
				if err != nil {
					retries -= 1
				} else {
					break
				}
			}
			return response.StatKey(), nil
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		requestQueue <- job
	}
	close(requestQueue)
	<-wp.Done
	b.StopTimer()
	report := collector.GenerateReport()
	b.Log(report)
}
