package benchmarks

import (
	"context"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
)

func CreateJob(client *HTTPClient) wpool.Job {
	return wpool.Job{
		ExecFn: func(ctx context.Context) (string, int64, error) {
			var (
				err      error
				response ResponseDetails
				target   pb.UserAccount
				retries  int64
			)
			for {
				response, err = client.MakeRequest(ctx, "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c", &target)
				if err == nil {
					break
				}
				retries++
			}
			return response.StatKey(), retries, err
		},
	}
}
func BenchmarkHTTP11GetWithWokers(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// start N workers
	wp := wpool.New(1)
	go wp.Run(ctx)

	requestQueue := wp.JobQueue()

	// collect stats
	collector := wpool.NewCollector(wp.Results())
	go collector.Run(ctx)

	// create job
	client := NewHTTPClient(HTTP)
	job := CreateJob(client)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		requestQueue <- job
	}
	close(requestQueue)
	<-wp.Done
	b.StopTimer()
	collector.GenerateReport(b)
}

func BenchmarkHTTP11GetWithWokersSUb(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	cases := []struct {
		name     string
		workers  int
		protocol TransportProtocolVer
	}{
		{
			name:     "HTTP1GetUserById-worker1",
			workers:  1,
			protocol: HTTP,
		},
		{
			name:     "HTTP1GetUserById-worker2",
			workers:  2,
			protocol: HTTP,
		},
		{
			name:     "HTTP1GetUserById-worker4",
			workers:  4,
			protocol: HTTP,
		},
		{
			name:     "HTTP1GetUserById-worker8",
			workers:  8,
			protocol: HTTP,
		},
		{
			name:     "HTTP1GetUserById-worker16",
			workers:  16,
			protocol: HTTP,
		},
		{
			name:     "HTTP1GetUserById-worker32",
			workers:  32,
			protocol: HTTP,
		},
		{
			name:     "HTTP1GetUserById-worker64",
			workers:  64,
			protocol: HTTP,
		},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			// start N workers
			wp := wpool.New(c.workers)
			go wp.Run(ctx)

			requestQueue := wp.JobQueue()

			// collect stats
			collector := wpool.NewCollector(wp.Results())
			go collector.Run(ctx)

			// create job
			client := NewHTTPClient(c.protocol)
			job := CreateJob(client)

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				requestQueue <- job
			}
			close(requestQueue)
			<-wp.Done
			b.StopTimer()
			collector.GenerateReport(b)
		})
	}
}
