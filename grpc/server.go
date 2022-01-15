package main

import (
	"context"
	"log"
	"net"

	"github.com/TheMickeyMike/grpc-rest-bench/data"
	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type API struct {
	server *grpc.Server
}

func NewAPI(userService pb.UsersServer) *API {
	var opts []grpc.ServerOption
	certFile := data.Path("x509/server.crt")
	keyFile := data.Path("x509/server.key")
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}

	server := grpc.NewServer(opts...)
	pb.RegisterUsersServer(server, userService)
	return &API{server: server}
}

func (api *API) Start(lst net.Listener) chan error {
	runtimeError := make(chan error)
	logger.Info("api : grpc api listening", zap.String("server", lst.Addr().String()))
	go func() {
		runtimeError <- api.server.Serve(lst)
	}()
	return runtimeError
}

func (api *API) Stop(ctx context.Context) error {
	logger.Info("api : start shutdown...")
	defer logger.Info("api : shutdown completed")

	// Asking listener to shutdown and load shed.
	stopped := make(chan struct{})
	go func() {
		api.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		api.server.Stop()
		logger.Info("api : graceful shutdown did not complete in time", zap.Error(ctx.Err()))
		return ctx.Err()
	case <-stopped:
		return nil
	}
}
