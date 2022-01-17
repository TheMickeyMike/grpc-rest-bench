package benchmarks

import (
	"context"
	"strings"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
)

func BenchmarkHTTP2GetWithWokers(b *testing.B) {
	client := NewHTTPClient(HTTP2)

	wp := wpool.New(1)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx) //start workers

	requestQueue := wp.JobQueue()

	collector := wpool.NewCollector(wp.Results())

	go collector.Run(ctx) //start result collector

	job := &wpool.Job{
		ExecFn: func(ctx context.Context, args interface{}) ([]string, error) {
			var respBody pb.UserAccount
			reqStats, err := client.MakeRequest(ctx, "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c", &respBody)
			if err != nil {
				return nil, err
			}
			// return []string{"200_OK", "HTTP_1_1"}, nil
			return []string{transformToStatKey(reqStats)}, nil
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

func transformToStatKey(stats []string) string {
	key := strings.Join(stats, "_")
	return strings.ReplaceAll(key, " ", "_")
}
