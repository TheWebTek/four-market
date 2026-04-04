package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap/zapcore"
)

// dailyFileWriter implements zapcore.WriteSyncer for daily log file rotation.
// This follows the Single Responsibility Principle by handling only file writing logic.
// It automatically creates a new log file when the date changes, enabling daily rotation.
//
// The writer implements the Liskov Substitution Principle by implementing the
// zapcore.WriteSyncer interface, allowing it to be used anywhere a WriteSyncer is expected.
type dailyFileWriter struct {
	basePath    string   // Directory where log files are stored
	filePrefix  string   // Prefix for log file names
	file        *os.File // Current open file handle
	currentDate string   // Currently tracked date for rotation
}

// newDailyFileWriter creates a new dailyFileWriter instance.
// It returns a zapcore.WriteSyncer interface, following the Interface Segregation principle.
// The returned type implements the WriteSyncer interface which only requires Write() and Close() methods.
func newDailyFileWriter(logPath, filePrefix string) (zapcore.WriteSyncer, error) {
	w := &dailyFileWriter{
		basePath:   logPath,
		filePrefix: filePrefix,
	}
	return zapcore.AddSync(w), nil
}

// Write writes the given bytes to the log file.
// It checks the current date on each write and rotates to a new file if the date has changed.
// This ensures that each day has its own log file (e.g., four-market-2026-04-04.log).
//
// Error handling:
// - Returns error if file cannot be opened/created
// - Logs error to stdout if file close fails (best effort, non-fatal)
// - Returns any write errors to the caller
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
