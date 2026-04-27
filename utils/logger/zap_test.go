package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	tmpDir := t.TempDir()
	log := newLogger(WithLogPath(tmpDir), WithFilePrefix("test"))
	assert.NotNil(t, log)
}

func TestNewLogger_WritesToFile(t *testing.T) {
	tmpDir := t.TempDir()
	log := newLogger(WithLogPath(tmpDir), WithFilePrefix("test"))
	log.Info("test message")
	log.Sync()
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Greater(t, len(files), 0)
}

func TestLogger_AllMethods(t *testing.T) {
	tmpDir := t.TempDir()
	log := newLogger(WithLogPath(tmpDir), WithFilePrefix("test"))
	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
}

func TestLogFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	log := newLogger(WithLogPath(tmpDir), WithFilePrefix("content"))
	log.Info("unique-message-12345")
	log.Sync()
	expectedFile := filepath.Join(tmpDir, "content-"+time.Now().Format("2006-01-02")+".log")
	content, err := os.ReadFile(expectedFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "unique-message-12345")
}
