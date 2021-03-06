package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheMickeyMike/grpc-rest-bench/user"
	"go.uber.org/zap"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func main() {
	// app main context
	ctx := context.Background()

	if err := run(ctx); err != nil {
		logger.Error("error :", zap.Error(err))
		os.Exit(1)
	}
	os.Exit(0)
}

func run(ctx context.Context) error {

	logger.Info("Starting....")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prepare listener for gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		return fmt.Errorf("main : failed to listen: %w", err)
	}

	// create api server
	grpcServer := NewAPI(user.NewService((user.NewDb(logger))))

	// server := NewAPI(":8080", time.Second*2, time.Second*2, router)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, shutdownSignals...)
	defer func() {
		signal.Stop(shutdown)
		close(shutdown)
	}()

	serverErrorsCh := make(chan error, 1)
	defer close(serverErrorsCh)

	// Start API server

	// Blocking main and waiting for shutdown.
	select {
	case err := <-grpcServer.Start(lis):
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Info("main : Start shutdown", zap.String("signal", sig.String()))

		// Asking listener to shutdown
		context, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := grpcServer.Stop(context)

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	case <-ctx.Done():
		logger.Info("main : context done")
	}
	return nil
}
