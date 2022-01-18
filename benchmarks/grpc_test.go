package benchmarks

import (
	"context"
	"log"
	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/data"
	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"github.com/TheMickeyMike/grpc-rest-bench/wpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func BenchmarkGRPCHTTP2GetWithWokers(b *testing.B) {
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/server.crt"), "example.com")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewUsersClient(conn)

	wp := wpool.New(1)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx) //start workers

	requestQueue := wp.JobQueue()

	collector := wpool.NewCollector(wp.Results())

	go collector.Run(ctx) //start result collector

	job := wpool.Job{
		ExecFn: func(ctx context.Context) (string, int64, error) {
			_, err := client.GetUser(ctx, &pb.UserRequest{Id: "61df07d341ed08ad981c143c"})
			if err != nil {
				return "ERROR", 0, nil
			}
			return "OK", 0, nil
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		requestQueue <- job
	}
	close(requestQueue)
	<-wp.Done
	b.StopTimer()
	collector.GenerateReport(b)
	// b.Log(report)
}
