package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

type ZapSugaredLogger = zap.SugaredLogger

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

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

type Config struct {
	LogPath    string
	FilePrefix string
	MinLevel   Level
	EncodeJSON bool
}

type Option func(*Config)

func WithLogPath(path string) Option {
	return func(c *Config) {
		c.LogPath = path
	}
}

func WithFilePrefix(prefix string) Option {
	return func(c *Config) {
		c.FilePrefix = prefix
	}
}

func WithMinLevel(level Level) Option {
	return func(c *Config) {
		c.MinLevel = level
	}
}

func WithJSONEncoding(encodeJSON bool) Option {
	return func(c *Config) {
		c.EncodeJSON = encodeJSON
	}
}
