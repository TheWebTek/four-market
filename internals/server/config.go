package server

import (
	"time"
)

type Config struct {
	Port            string
	GracefulTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

type Option func(*Config)

func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithGracefulTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.GracefulTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.WriteTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.IdleTimeout = timeout
	}
}

func defaultConfig() Config {
	return Config{
		Port:            "8080",
		GracefulTimeout: 10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
	}
}
