package benchmarks

import (
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/data"
	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GenerateBenchmarkGrpcCases() []BenchmarkCase {
	var (
		cases        []BenchmarkCase
		workersCount = []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
	)
	for _, worker := range workersCount {
		cases = append(cases, BenchmarkCase{
			name:    strconv.Itoa(worker),
			workers: worker,
		})
	}
	return cases
}

func CreateGrpcAPIJob(client pb.UsersClient) wpool.Job {
	return wpool.Job{
		ExecFn: func(ctx context.Context) (string, int64, error) {
			_, err := client.GetUser(ctx, &pb.UserRequest{Id: "61df07d341ed08ad981c143c"})
			if err != nil {
				return "ERROR", 0, nil
			}
			return "OK", 0, nil
		},
	}
}

func BenchmarkGrpcGetUserById(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	cases := GenerateBenchmarkGrpcCases()

	//setup grpc client
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/server.crt"), "example.com")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewUsersClient(conn)

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
			job := CreateGrpcAPIJob(client)

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
