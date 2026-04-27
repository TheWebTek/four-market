package logger

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// singleton struct holds the thread-safe singleton instance.
type singleton struct {
	logger *zap.SugaredLogger
	once   sync.Once
}

var instance singleton

var defaultConfig = Config{
	LogPath:    "logs",
	FilePrefix: "four-market",
	MinLevel:   zapcore.InfoLevel,
	EncodeJSON: true,
}

func (s *singleton) getInstance(opts ...Option) Logger {
	s.once.Do(func() {
		s.logger = newLogger(opts...)
	})
	return s.logger
}

func newLogger(opts ...Option) *zap.SugaredLogger {
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
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
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

	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig()),
			writer,
			cfg.MinLevel,
		),
		zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			cfg.MinLevel,
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger.Sugar()
}

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

// GetInstance returns the singleton logger instance.
// It uses sync.Once to ensure thread-safe lazy initialization.
// Only the first call creates the logger; subsequent calls return the same instance.
//
// Options can be passed on the first call to configure the logger.
// These options are ignored on subsequent calls.
func GetInstance(opts ...Option) Logger {
	return instance.getInstance(opts...)
}

// Sync flushes any buffered log entries.
// It should be called during application shutdown to ensure all logs are written.
func Sync() error {
	return instance.logger.Sync()
}
