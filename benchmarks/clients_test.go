package benchmarks

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"testing"

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
	"golang.org/x/net/http2"
)

const wcount = 6

var client http.Client

func init() {
	client = http.Client{}
}

type Request struct {
	Path string
}

type Result struct {
	User       warehouse.UserAccount
	StatusCode int
	Error      error
}

func createTLSConfigWithCustomCert() *tls.Config {
	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCert, err := ioutil.ReadFile("../rest/ssl/server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("../rest/ssl/server.crt", "../rest/ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}

	// Create TLS configuration with the certificate of the server
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
	}
}

func MakeRequest(ctx context.Context, url string, output interface{}) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// decode input
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(output)

	return res.StatusCode, nil
}

func BenchmarkRestHTTP2GetWithWokers(b *testing.B) {
	var wg sync.WaitGroup

	client.Transport = &http2.Transport{
		AllowHTTP:       true,
		TLSClientConfig: createTLSConfigWithCustomCert(),
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	requestQueue := make(chan Request, wcount)
	resultsQueue := make(chan Result, wcount)

	wg.Add(wcount)
	for i := 0; i < wcount; i++ {
		go worker(ctx, &wg, requestQueue, resultsQueue, i)
	}

	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		requestQueue <- Request{Path: "http://localhost:8080/api/v1/users/61df07d341ed08ad981c143c"}
	}
	close(requestQueue)
	wg.Wait()
	close(resultsQueue)
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
		go worker(ctx, &wg, requestQueue, resultsQueue, i)
	}

	b.ResetTimer() // don't count worker initialization time
	for n := 0; n < b.N; n++ {
		requestQueue <- Request{Path: "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c"}
	}
	close(requestQueue)
	wg.Wait()
	close(resultsQueue)
}

func worker(ctx context.Context, wg *sync.WaitGroup, requestQueue <-chan Request, resultsQueue chan<- Result, id int) {
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
			resultsQueue <- Result{
				Error: ctx.Err(),
			}
			return
		}
	}
}
