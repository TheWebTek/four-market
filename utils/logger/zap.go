package logger

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// defaultLogger is the singleton logger instance used throughout the application.
// It is initialized during package init() and provides a convenient global logging instance.
var defaultLogger *zap.SugaredLogger

// defaultConfig defines the default configuration values for the logger.
// These values are used when no custom options are provided.
var defaultConfig = Config{
	LogPath:    "logs",
	FilePrefix: "four-market",
	MinLevel:   zapcore.InfoLevel,
	EncodeJSON: true,
}

// init initializes the default logger when the package is first loaded.
// This ensures logging is ready before any other code runs.
// It follows the Single Responsibility Principle by handling only initialization.
func init() {
	defaultLogger = New()
}

// New creates and configures a new zap.SugaredLogger instance.
// It applies the provided functional options to customize the logger configuration.
// If no options are provided, default configuration values are used.
// Environment variables (LOG_PATH, LOG_FILE_PREFIX) can override default values.
//
// The function follows the Dependency Inversion principle by returning an interface
// (specifically *zap.SugaredLogger which implements Logger), allowing consumers
// to depend on the abstraction rather than the concrete implementation.
//
// Example usage:
//
//	logger := logger.New(
//		logger.WithLogPath("/var/log/myapp"),
//		logger.WithFilePrefix("myapp"),
//		logger.WithMinLevel(logger.DebugLevel),
//	)
func New(opts ...Option) *zap.SugaredLogger {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	_ = godotenv.Load()
	if cfg.LogPath == "" {
		cfg.LogPath = os.Getenv("LOG_PATH")
	}
	if cfg.LogPath == "" {
		cfg.LogPath = "logs"
	}
	if cfg.FilePrefix == "" {
		cfg.FilePrefix = os.Getenv("LOG_FILE_PREFIX")
	}
	if cfg.FilePrefix == "" {
		cfg.FilePrefix = "four-market"
	}

	if err := os.MkdirAll(cfg.LogPath, 0o755); err != nil {
		println("Failed to create log directory: " + err.Error())
	}

	writer, err := newDailyFileWriter(cfg.LogPath, cfg.FilePrefix)
	if err != nil {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig()),
			zapcore.AddSync(os.Stdout),
			cfg.MinLevel,
		)
		return zap.New(core, zap.AddCaller()).Sugar()
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig()),
		writer,
		cfg.MinLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger.Sugar()
}

// encoderConfig returns the default encoder configuration for structured logging.
// It defines how log fields are encoded, including timestamp format, level encoding,
// duration format, and caller information.
func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// GetInstance returns the default logger singleton instance.
// This provides convenient access to the pre-configured logger without
// requiring manual initialization.
func GetInstance() *zap.SugaredLogger {
	return defaultLogger
}

// Sync flushes any buffered log entries.
// It should be called during application shutdown to ensure all logs are written.
// Returns an error if the sync operation fails.
func Sync() error {
	return defaultLogger.Sync()
}
