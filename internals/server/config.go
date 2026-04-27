// Package server provides HTTP server functionality with Echo framework.
package server

import "time"

// Config holds the configuration options for the HTTP server.
// Following the Open/Closed principle, new configuration options can be added
// without modifying the core server implementation.
type Config struct {
	Port            string        // Server listen address (e.g., "8080" or ":8080")
	GracefulTimeout time.Duration // Maximum time to wait for pending requests during shutdown
	ReadTimeout     time.Duration // Maximum time to read the full request (including body)
	WriteTimeout    time.Duration // Maximum time to write the response
	IdleTimeout     time.Duration // Maximum time to wait for the next request when keep-alive is enabled
}

// Option is a functional option pattern for configuring the server.
// This follows the Open/Closed principle, allowing flexible configuration
// without modifying the core server implementation.
type Option func(*Config)

// WithPort sets the server port.
// Default: "8080"
func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithGracefulTimeout sets the graceful shutdown timeout.
// This is the maximum time the server will wait for in-flight requests to complete
// before force closing connections during shutdown.
// Default: 10 seconds
func WithGracefulTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.GracefulTimeout = timeout
	}
}

// WithReadTimeout sets the maximum time allowed to read the entire request,
// including the body. This helps prevent slow clients from exhausting server resources.
// Default: 30 seconds
func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets the maximum time allowed to write the response.
// This helps prevent slow clients from keeping connections open too long.
// Default: 30 seconds
func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.WriteTimeout = timeout
	}
}

// WithIdleTimeout sets the maximum time to wait for the next request when
// keep-alive is enabled. If no new request arrives within this time,
// the connection will be closed.
// Default: 60 seconds
func WithIdleTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.IdleTimeout = timeout
	}
}

// WithLogger sets a custom logger instance.
// If not provided, the default singleton logger will be used.
func WithLogger(log logger.Logger) Option {
	return func(c *Config) {
		c.Logger = log
	}
}

// defaultConfig returns the default configuration values.
// These values are used when no custom options are provided.
func defaultConfig() Config {
	return Config{
		GracefulTimeout: 10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
	}
}
