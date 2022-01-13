package benchmarks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
)

const wcount = 3

type Request struct {
	Path string
}

type Result struct {
	User       warehouse.UserAccount
	StatusCode int
	Error      error
}

func MakeRequest(ctx context.Context, url string, output interface{}) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// decode input
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(output)

	return res.StatusCode, nil
}

func BenchmarkRest_10MB(b *testing.B) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	requestQueue := make(chan Request, wcount)
	results := make(chan Result, wcount)

	wg.Add(wcount)
	for i := 0; i < wcount; i++ {
		go worker(ctx, &wg, requestQueue, results, i)
	}

	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		requestQueue <- Request{Path: "http://localhost:8080/api/v1/users/61df07d341ed08ad981c143c"}
	}
	close(requestQueue)
	wg.Wait()
	// for s := range results {
	// 	log.Printf("RESULT: %+v\n", s)
	// }
	close(results)
}

func worker(ctx context.Context, wg *sync.WaitGroup, requestQueue <-chan Request, results chan<- Result, id int) {
	// log.Printf("[worker %d]start\n", id)
	defer wg.Done()
	for {
		select {
		case req, ok := <-requestQueue:
			if !ok {
				return
			}
			var users []warehouse.UserAccount
			_, err := MakeRequest(ctx, req.Path, &users)
			if err != nil {
				log.Printf("[worker %d] die (reason: %s)\n", id, err)
				return
			}
			// respCode, err := MakeRequest(ctx, req.Path, &users)
			// log.Printf("[worker %d]  result: %+v code: %d, err: %v\n", id, users, respCode, err)
			// results <- Result{user, respCode, err}
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Error: ctx.Err(),
			}
			return
		}
	}
}
