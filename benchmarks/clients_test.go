package benchmarks

import (
	"context"
	"net/http"
	"sync"

	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
	"golang.org/x/net/http2"
)

const wcount = 128

func BenchmarkRestHTTP2GetWithWokers(b *testing.B) {
	var (
		wWg sync.WaitGroup
		cWg sync.WaitGroup
	)

	client.Transport = &http2.Transport{
		TLSClientConfig: createTLSConfigWithCustomCert(),
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	collector := Collector{make(map[string]int)}

	requestQueue := make(chan *Request, wcount)
	resultsQueue := make(chan *Result, wcount*15) // make channel bigger to not block writing result from workers

	collector.Start(&cWg, resultsQueue)

	wWg.Add(wcount)
	for i := 0; i < wcount; i++ {
		go Worker(ctx, &wWg, requestQueue, resultsQueue, i)
	}

	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		requestQueue <- &Request{Path: "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c", ResponseObject: &warehouse.UserAccount{}}
		// requestQueue <- &Request{Path: "https://localhost:8080/api/v1/small", ResponseObject: &warehouse.SmallResponse{}}
	}
	close(requestQueue) //stop workers
	wWg.Wait()          //wait for workers gracefull shutdown
	close(resultsQueue) //close result channel and collector
	b.StopTimer()       //stop benchmark timer, all request has been made
	cWg.Wait()          // wait for collecting result
	b.Logf("%+v SUM: %10d", collector.GetResults(), collector.GetSum())
}

func BenchmarkRestHTTP11Get(b *testing.B) {
	var (
		wWg sync.WaitGroup
		cWg sync.WaitGroup
	)

	client.Transport = &http.Transport{
		TLSClientConfig: createTLSConfigWithCustomCert(),
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	collector := Collector{make(map[string]int)}

	requestQueue := make(chan *Request, wcount)
	resultsQueue := make(chan *Result, wcount*15) // make channel bigger to not block writing result from workers

	collector.Start(&cWg, resultsQueue)

	wWg.Add(wcount)
	for i := 0; i < wcount; i++ {
		go Worker(ctx, &wWg, requestQueue, resultsQueue, i)
	}

	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		requestQueue <- &Request{Path: "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c", ResponseObject: &warehouse.UserAccount{}}
		// requestQueue <- &Request{Path: "https://localhost:8080/api/v1/small", ResponseObject: &warehouse.SmallResponse{}}
	}
	close(requestQueue) //stop workers
	wWg.Wait()          //wait for workers gracefull shutdown
	close(resultsQueue) //close result channel and collector
	b.StopTimer()       //stop benchmark timer, all request has been made
	cWg.Wait()          // wait for collecting result
	b.Logf("%+v SUM: %10d", collector.GetResults(), collector.GetSum())
}
