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

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
)

var client http.Client

func init() {
	client = http.Client{}
}

type Result struct {
	User       warehouse.UserAccount
	StatusCode int
	Error      error
}

type Request struct {
	Path           string
	ResponseObject *warehouse.SmallResponse
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

	// use discard when not reading body
	// io.Copy(ioutil.Discard, res.Body)

	// use decode when
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(output); err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}

func Worker(ctx context.Context, wg *sync.WaitGroup, requestQueue <-chan Request, resultsQueue chan<- Result, id int) {
	defer wg.Done()
	for {
		select {
		case req, ok := <-requestQueue:
			if !ok {
				return
			}
			_, err := MakeRequest(ctx, req.Path, req.ResponseObject)
			if err != nil {
				log.Printf("[worker %d] die (reason: %s)\n", id, err)
				return
			}
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			resultsQueue <- Result{
				Error: ctx.Err(),
			}
			return
		}
	}
}
