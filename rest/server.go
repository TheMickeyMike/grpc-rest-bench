package main

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/TheMickeyMike/grpc-rest-bench/data"
	"go.uber.org/zap"
)

type API struct {
	server http.Server
}

func NewAPI(host string, readTimeout, writeTimeout time.Duration, router http.Handler) *API {
	return &API{
		http.Server{
			Addr:           host,
			Handler:        router,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      tlsConfig(),
		},
	}
}

func (api *API) Start(serverErrors chan<- error) {
	go func() {
		logger.Info("api : API Listening", zap.String("server", api.server.Addr))
		err := api.server.ListenAndServeTLS("", "")
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()
}

func (api *API) Stop(ctx context.Context) error {
	logger.Info("api : Start shutdown...")
	defer logger.Info("api : Completed")

	// Asking listener to shutdown and load shed.
	if err := api.server.Shutdown(ctx); err != nil {
		logger.Info("api : Graceful shutdown did not complete in time", zap.Error(err))
		if err := api.server.Close(); err != nil {
			return err
		}
	}
	return nil
}

func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile(data.Path("x509/server.crt"))
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(data.Path("x509/server.key"))
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
	}
}
