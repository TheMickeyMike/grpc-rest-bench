package benchmarks

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/TheMickeyMike/grpc-rest-bench/data"
	"golang.org/x/net/http2"
)

type TransportProtocolVer string

var (
	HTTP  TransportProtocolVer = "HTTP1"
	HTTP2 TransportProtocolVer = "HTTP2"
)

type ResponseDetails struct {
	Status string
	Proto  string
}

func (r *ResponseDetails) StatKey() string {
	stats := []string{r.Status, r.Proto}
	key := strings.Join(stats, "_")
	return strings.ReplaceAll(key, " ", "_")
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(protocol TransportProtocolVer) *HTTPClient {
	c := &http.Client{}
	switch protocol {
	case HTTP:
		c.Transport = &http.Transport{
			TLSClientConfig: createTLSConfigWithCustomCert(),
		}

	case HTTP2:
		c.Transport = &http2.Transport{
			TLSClientConfig: createTLSConfigWithCustomCert(),
		}
	}
	return &HTTPClient{c}
}

func (c *HTTPClient) MakeRequest(ctx context.Context, url string, output interface{}) (ResponseDetails, error) {
	result := ResponseDetails{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		// log.Println(err) //socket: too many open files https://github.com/golang/go/issues/18588
		return ResponseDetails{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body) // drop resp: io.Copy(ioutil.Discard, res.Body)
	if err = decoder.Decode(output); err != nil {
		return result, err
	}

	result.Status = res.Status
	result.Proto = res.Proto

	return result, nil
}

func createTLSConfigWithCustomCert() *tls.Config {
	caCert, err := ioutil.ReadFile(data.Path("x509/server.crt"))
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server.crt"), data.Path("x509/server.key"))
	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
}
