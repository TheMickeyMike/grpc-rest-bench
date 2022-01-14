package benchmarks

import (
	"context"
	"net/http"
	"sync"

	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
	"golang.org/x/net/http2"
)

const wcount = 150

func BenchmarkRestHTTP2GetWithWokers(b *testing.B) {
	var wg sync.WaitGroup

	client.Transport = &http2.Transport{
		TLSClientConfig: createTLSConfigWithCustomCert(),
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	requestQueue := make(chan Request, wcount)
	resultsQueue := make(chan Result, wcount)

	for i := 0; i < wcount; i++ {
		go Worker(ctx, &wg, requestQueue, resultsQueue, i)
	}
	wg.Add(wcount)
	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		// requestQueue <- Request{Path: "http://localhost:8080/api/v1/users/61df07d341ed08ad981c143c"}
		requestQueue <- Request{Path: "https://127.0.0.1:8080/api/v1/small", ResponseObject: &warehouse.SmallResponse{}}
	}
	b.StopTimer()       //stop benchmark timer
	close(requestQueue) //stop workers
	wg.Wait()           //wait for workers gracefull shutdown
	close(resultsQueue) //close result channel
}

func BenchmarkRestHTTP11Get(b *testing.B) {
	var wg sync.WaitGroup

	client.Transport = &http.Transport{
		TLSClientConfig: createTLSConfigWithCustomCert(),
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	requestQueue := make(chan Request, wcount)
	resultsQueue := make(chan Result, wcount)

	wg.Add(wcount)
	for i := 0; i < wcount; i++ {
		go Worker(ctx, &wg, requestQueue, resultsQueue, i)
	}

	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		// requestQueue <- Request{Path: "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c"}
		requestQueue <- Request{Path: "https://127.0.0.1:8080/api/v1/small"}
	}
	close(requestQueue) //stop workers
	wg.Wait()           //wait for workers gracefull shutdown
	close(resultsQueue) //close result channel
}
