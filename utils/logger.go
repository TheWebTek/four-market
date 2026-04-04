package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = newLogger()

func newLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02 Jan 03:04:05 PM")

	logger, err := config.Build()
	if err != nil {
		fallbackConfig := zap.NewDevelopmentConfig()
		fallbackConfig.Encoding = "console"
		fallbackConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		fallbackConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("02 Jan 03:04:05 PM"))
		}

		logger, _ = fallbackConfig.Build()
	}

	return logger
}

func GetInstance() *zap.SugaredLogger {
	return log.Sugar()
}
