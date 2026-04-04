package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap/zapcore"
)

type dailyFileWriter struct {
	basePath    string
	filePrefix  string
	file        *os.File
	currentDate string
}

func newDailyFileWriter(logPath, filePrefix string) (zapcore.WriteSyncer, error) {
	w := &dailyFileWriter{
		basePath:   logPath,
		filePrefix: filePrefix,
	}
	return zapcore.AddSync(w), nil
}

func (w *dailyFileWriter) Write(p []byte) (n int, err error) {
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
