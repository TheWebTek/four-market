package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDailyFileWriter(t *testing.T) {
	tmpDir := t.TempDir()
	w, err := newDailyFileWriter(tmpDir, "test")
	require.NoError(t, err)
	assert.NotNil(t, w)
}

func TestDailyFileWriter_Write(t *testing.T) {
	tmpDir := t.TempDir()
	w, err := newDailyFileWriter(tmpDir, "test")
	require.NoError(t, err)
	writer := w.(interface{ Write([]byte) (int, error) })
	n, err := writer.Write([]byte("test\n"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0)
}

func TestDailyFileWriter_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	w, _ := newDailyFileWriter(tmpDir, "test")
	writer := w.(interface{ Write([]byte) (int, error) })
	writer.Write([]byte("msg\n"))
	expectedFile := filepath.Join(tmpDir, "test-"+time.Now().Format("2006-01-02")+".log")
	_, err := os.Stat(expectedFile)
	assert.NoError(t, err)
}
