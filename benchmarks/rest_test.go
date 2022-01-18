package benchmarks

import (
	"context"
	"strconv"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
)

type BenchmarkCase struct {
	name     string
	workers  int
	protocol TransportProtocolVer
}

func GenerateBenchmarkHTTPCases(protoVer TransportProtocolVer) []BenchmarkCase {
	var (
		cases        []BenchmarkCase
		workersCount = []int{1, 2, 4, 8, 16, 32, 64}
	)
	for _, worker := range workersCount {
		cases = append(cases, BenchmarkCase{
			name:     strconv.Itoa(worker),
			workers:  worker,
			protocol: protoVer,
		})
	}
	return cases
}

func CreateRestAPIJob(client *HTTPClient) wpool.Job {
	return wpool.Job{
		ExecFn: func(ctx context.Context) (string, int64, error) {
			var (
				err      error
				response ResponseDetails
				target   pb.UserAccount
				retries  int64
			)
			for {
				response, err = client.MakeRequest(ctx, "https://127.0.0.1:8080/api/v1/users/61df07d341ed08ad981c143c", &target)
				if err == nil {
					break
				}
				retries++
			}
			return response.StatKey(), retries, err
		},
	}
}

func BenchmarkHTTP11GetUserById(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	cases := GenerateBenchmarkHTTPCases(HTTP)

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
			job := CreateRestAPIJob(client)

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

func BenchmarkHTTP2GetUserById(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	cases := GenerateBenchmarkHTTPCases(HTTP2)

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
			job := CreateRestAPIJob(client)

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
