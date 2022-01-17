package benchmarks

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TheMickeyMike/grpc-rest-bench/data"
	"golang.org/x/net/http2"
)

type TransportProtocolVer string

var (
	HTTP  TransportProtocolVer = "HTTP1"
	HTTP2 TransportProtocolVer = "HTTP2"
)

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

func (c *HTTPClient) MakeRequest(ctx context.Context, url string, output interface{}) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		// log.Println(err) //socket: too many open files https://github.com/golang/go/issues/18588
		return nil, err
	}
	defer res.Body.Close()

	// use discard when not reading body
	// io.Copy(ioutil.Discard, res.Body)

	// use decode when
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(output); err != nil {
		return nil, err
	}
	return []string{res.Status, res.Proto}, nil
}

func createTLSConfigWithCustomCert() *tls.Config {
	// Create a pool with the server certificate since it is not signed
	// by a known CA
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

	// Create TLS configuration with the certificate of the server
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
}
