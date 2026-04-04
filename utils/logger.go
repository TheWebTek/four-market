package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// init initializes the logger when the package is first loaded.
// This ensures logging is ready before any other code runs.
func init() {
	log = newLogger()
}

// loadEnv reads the LOG_PATH and LOG_FILE_PREFIX environment variables.
// LOG_PATH: directory where log files will be stored (default: "logs")
// LOG_FILE_PREFIX: prefix for log file names (default: "four-market")
func loadEnv() (string, string) {
	_ = godotenv.Load()
	return os.Getenv("LOG_PATH"), os.Getenv("LOG_FILE_PREFIX")
}

// newLogger creates and configures a new zap logger instance.
// It sets up:
// - A custom daily rotating file writer
// - JSON encoding for structured logs
// - Caller information for trace debugging
func newLogger() *zap.Logger {
	logPath, filePrefix := loadEnv()

	if logPath == "" {
		logPath = "logs"
	}

	if filePrefix == "" {
		filePrefix = "four-market"
	}

	if err := os.MkdirAll(logPath, 0o755); err != nil {
		println("Failed to create log directory: " + err.Error())
	}

	dateStr := time.Now().Format("2006-01-02")
	_ = dateStr // used in dailyWriter for dynamic date checking

	writer := zapcore.AddSync(&dailyWriter{
		basePath:   logPath,
		filePrefix: filePrefix,
	})

	encoderConfig := zapcore.EncoderConfig{
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

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger
}

// dailyWriter implements zapcore.WriteSyncer for daily log file rotation.
// It automatically creates a new log file when the date changes.
type dailyWriter struct {
	basePath    string
	filePrefix  string
	file        *os.File
	currentDate string
}

// Write writes the given bytes to the log file.
// It checks the current date on each write and rotates to a new file if the date changed.
func (w *dailyWriter) Write(p []byte) (n int, err error) {
	currentDate := time.Now().Format("2006-01-02")
	expectedFileName := w.filePrefix + "-" + currentDate + ".log"

	if w.file == nil || w.currentDate != currentDate {
		if w.file != nil {
			if closeErr := w.file.Close(); closeErr != nil {
				println("Failed to close log file: " + closeErr.Error())
			}
		}

		w.currentDate = currentDate
		filePath := filepath.Join(w.basePath, expectedFileName)

		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			println("Failed to open log file: " + err.Error())
			return 0, err
		}
		w.file = file
	}

	return w.file.Write(p)
}

// GetInstance returns a SugaredLogger instance for logging.
// This provides a more convenient API with formatted logging support.
func GetInstance() *zap.SugaredLogger {
	return log.Sugar()
}
