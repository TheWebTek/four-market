// Package logger provides a structured logging solution with daily file rotation.
//
// This package implements the SOLID principles:
// - Single Responsibility: Each file handles one concern (interfaces, implementation, file writing)
// - Open/Closed: Extend functionality through Options without modifying core code
// - Liskov Substitution: Logger interface allows different implementations
// - Interface Segregation: Focused interfaces for specific needs
// - Dependency Inversion: Depend on abstractions (Logger interface) not concretions
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level represents the severity of log messages.
// It is an alias for zapcore.Level to maintain compatibility.
type Level = zapcore.Level

// ZapSugaredLogger is an alias for zap.SugaredLogger.
// It provides a more convenient API with formatted logging support.
type ZapSugaredLogger = zap.SugaredLogger

const (
	DebugLevel = zapcore.DebugLevel // Debug level for detailed debugging information
	InfoLevel  = zapcore.InfoLevel  // Info level for general informational messages
	WarnLevel  = zapcore.WarnLevel  // Warn level for warning messages
	ErrorLevel = zapcore.ErrorLevel // Error level for error messages
	FatalLevel = zapcore.FatalLevel // Fatal level for critical errors that cause termination
)

// Logger interface defines the contract for logging operations.
// This interface follows the Interface Segregation principle from SOLID,
// allowing consumers to depend only on the logging methods they need.
type Logger interface {
	Debug(args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync() error
}

// Config holds the configuration options for the logger.
// Following the Open/Closed principle, new configuration options can be added
// without modifying the existing logger implementation.
type Config struct {
	LogPath    string // Directory where log files will be stored (default: "logs")
	FilePrefix string // Prefix for log file names (default: "four-market")
	MinLevel   Level  // Minimum log level to output (default: InfoLevel)
	EncodeJSON bool   // Whether to encode logs as JSON (default: true)
}

// Option is a functional option pattern for configuring the logger.
// This follows the Open/Closed principle, allowing flexible configuration
// without modifying the core logger implementation.
type Option func(*Config)

// WithLogPath sets the directory where log files will be stored.
// Default: "logs"
func WithLogPath(path string) Option {
	return func(c *Config) {
		c.LogPath = path
	}
}

// WithFilePrefix sets the prefix for log file names.
// Default: "four-market"
func WithFilePrefix(prefix string) Option {
	return func(c *Config) {
		c.FilePrefix = prefix
	}
}

// WithMinLevel sets the minimum log level to output.
// Default: InfoLevel
func WithMinLevel(level Level) Option {
	return func(c *Config) {
		c.MinLevel = level
	}
}

// WithJSONEncoding sets whether to encode logs as JSON.
// Default: true
func WithJSONEncoding(encodeJSON bool) Option {
	return func(c *Config) {
		c.EncodeJSON = encodeJSON
	}
}
