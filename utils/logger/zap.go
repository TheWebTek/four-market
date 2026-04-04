package logger

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.SugaredLogger
var defaultConfig = Config{
	LogPath:    "logs",
	FilePrefix: "four-market",
	MinLevel:   zapcore.InfoLevel,
	EncodeJSON: true,
}

func init() {
	defaultLogger = New()
}

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

func GetInstance() *zap.SugaredLogger {
	return defaultLogger
}

func Sync() error {
	return defaultLogger.Sync()
}
