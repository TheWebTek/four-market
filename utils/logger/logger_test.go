package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithLogPath(t *testing.T) {
	cfg := &Config{}
	WithLogPath("logs")(cfg)
	assert.Equal(t, "logs", cfg.LogPath)
}

func TestWithFilePrefix(t *testing.T) {
	cfg := &Config{}
	WithFilePrefix("myapp")(cfg)
	assert.Equal(t, "myapp", cfg.FilePrefix)
}

func TestWithMinLevel(t *testing.T) {
	cfg := &Config{}
	WithMinLevel(DebugLevel)(cfg)
	assert.Equal(t, DebugLevel, cfg.MinLevel)
}

func TestWithJSONEncoding(t *testing.T) {
	cfg := &Config{}
	WithJSONEncoding(false)(cfg)
	assert.False(t, cfg.EncodeJSON)
}

func TestLoggerInterface(t *testing.T) {
	var _ Logger = (*zap.SugaredLogger)(nil)
}
