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
	"strings"
	"sync"
	"time"

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
)

var client http.Client

func init() {
	client = http.Client{}
}

type Result struct {
	User       warehouse.UserAccount
	StatusCode string
	Proto      string
	Error      error
	Retries    int
}

// type Request struct {
// 	Path           string
// 	ResponseObject *warehouse.SmallResponse
// }

type Request struct {
	Path           string
	ResponseObject *warehouse.UserAccount
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
	}
}

func MakeRequest(ctx context.Context, url string, output interface{}) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	res, err := client.Do(req)
	if err != nil {
		// log.Println(err) //socket: too many open files https://github.com/golang/go/issues/18588
		return "", "", err
	}
	defer res.Body.Close()

	// use discard when not reading body
	// io.Copy(ioutil.Discard, res.Body)

	// use decode when
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(output); err != nil {
		return "", "", err
	}
	return res.Status, res.Proto, nil
}

func Worker(ctx context.Context, wg *sync.WaitGroup, requestQueue <-chan *Request, resultsQueue chan<- *Result, id int) {
	defer wg.Done()
	for {
		select {
		case req, ok := <-requestQueue:
			if !ok {
				return
			}
			var (
				code, proto string
				err         error
			)
			code, proto, err = MakeRequest(ctx, req.Path, req.ResponseObject)
			var retry int
			if err != nil {
				for err != nil && retry < 3 {
					// fmt.Println("retry")
					time.Sleep(time.Second * 1)
					code, proto, err = MakeRequest(ctx, req.Path, req.ResponseObject)
					retry++
				}
			}

			resultsQueue <- &Result{StatusCode: code, Proto: proto, Error: err, Retries: retry}
			// if err != nil {

			// 	log.Printf("[worker %d] die (reason: %s)\n", id, err)
			// 	return
			// }
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			resultsQueue <- &Result{
				Error: ctx.Err(),
			}
			return
		}
	}
}

type Collector struct {
	results map[string]int
}

func (c *Collector) Start(wg *sync.WaitGroup, resultsQueue <-chan *Result) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for res := range resultsQueue {
			if res.Error != nil {
				c.results["error"] += 1
			} else if res.Retries > 0 {
				c.results["retries"] += 1
			} else {
				c.results[normalizeKeys(res.StatusCode, res.Proto)] += 1
			}
		}
	}()
}

func (c *Collector) GetResults() map[string]int {
	return c.results
}

func (c *Collector) GetSum() int {
	var sum int
	for _, v := range c.results {
		sum += v
	}
	return sum
}

func normalizeKeys(keys ...string) string {
	var resultKey string
	for _, key := range keys {
		resultKey += "_" + strings.ReplaceAll(key, " ", "_")
	}
	return resultKey
}
