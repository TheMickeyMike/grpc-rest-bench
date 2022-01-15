package main

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()

	// Override zap default logger
	zap.ReplaceGlobals(logger)
}
