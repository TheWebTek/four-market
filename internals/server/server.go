package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheWebTek/four-market/utils/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Server interface defines the contract for HTTP server implementations.
// This follows the Interface Segregation and Dependency Inversion principles,
// allowing consumers to depend on the abstraction (interface) rather than
// the concrete implementation (echoServer).
type Server interface {
	Start() error
}

// echoServer is the concrete implementation of the Server interface using Echo framework.
// It follows the Single Responsibility Principle by handling only server-related logic.
type echoServer struct {
	echo   *echo.Echo               // Echo engine instance
	config Config                   // Server configuration
	log    *logger.ZapSugaredLogger // Logger instance
}

// New creates a new Server instance with the provided options.
// It applies functional options to customize the server configuration.
// Environment variable PORT can be used to override the default port.
//
// Example usage:
//
//	srv := server.New(
//		server.WithPort("3000"),
//		server.WithGracefulTimeout(30 * time.Second),
//	)
func New(opts ...Option) Server {
	_ = godotenv.Load()

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Port == "" {
		cfg.Port = os.Getenv("PORT")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())

	return &echoServer{
		echo:   e,
		config: cfg,
		log:    logger.GetInstance(),
	}
}

// Start starts the HTTP server with graceful shutdown handling.
// It listens for OS interrupt signals (SIGINT, SIGTERM) and gracefully shuts down
// the server, allowing in-flight requests to complete within the configured timeout.
//
// The method:
// 1. Registers signal handlers for graceful shutdown
// 2. Starts the server with the configured port and graceful timeout
// 3. Blocks until a shutdown signal is received
// 4. Initiates graceful shutdown with the configured timeout
//
// Returns error if server fails to start (not on graceful shutdown).
func (s *echoServer) Start() error {
	s.log.Infow("starting server", "port", s.config.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startConfig := echo.StartConfig{
		Address:         ":" + s.config.Port,
		GracefulTimeout: s.config.GracefulTimeout,
		HideBanner:      true,
		HidePort:        true,
	}

	if err := startConfig.Start(ctx, s.echo); err != nil {
		s.log.Infow("server stopped", "error", err)
		return err
	}

	s.log.Info("server stopped")
	return nil
}

// Start is a convenience function that creates a new server with default configuration
// and starts it immediately. It uses the package-level default configuration.
func Start() {
	New().Start()
}
