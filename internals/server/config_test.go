package server

import (
	"testing"
	"time"

	"github.com/TheWebTek/four-market/utils/logger"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	assert.Equal(t, 10*time.Second, cfg.GracefulTimeout)
	assert.Empty(t, cfg.Port)
}

func TestWithPort(t *testing.T) {
	cfg := &Config{}
	WithPort("8080")(cfg)
	assert.Equal(t, "8080", cfg.Port)
}

func TestWithGracefulTimeout(t *testing.T) {
	cfg := &Config{}
	WithGracefulTimeout(20*time.Second)(cfg)
	assert.Equal(t, 20*time.Second, cfg.GracefulTimeout)
}

func TestWithLogger(t *testing.T) {
	cfg := &Config{}
	l := logger.GetInstance()
	WithLogger(l)(cfg)
	assert.Same(t, l, cfg.Logger)
}
