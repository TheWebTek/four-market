// Package main is the entry point for the four-market application.
//
// This package demonstrates the Dependency Inversion principle by depending on
// abstractions (logger.Logger interface via GetInstance, server.Server interface)
// rather than concrete implementations.
package main

import (
	"github.com/TheWebTek/four-market/internals/server"
	"github.com/TheWebTek/four-market/utils/logger"
)

func main() {
	// Get the default logger instance (singleton pattern)
	// The logger is pre-configured with daily file rotation
	log := logger.GetInstance()

	// Log server startup
	log.Info("Server Starting...")

	// Start the HTTP server with graceful shutdown
	// The server will listen for SIGINT/SIGTERM and gracefully shutdown
	server.Start()
}
