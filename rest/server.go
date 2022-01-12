package main

import (
	"context"
	"errors"
	"net/http"
	"time"

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
		},
	}
}

func (api *API) Start(serverErrors chan<- error) {
	go func() {
		logger.Info("api : API Listening", zap.String("server", api.server.Addr))
		if err := api.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
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
