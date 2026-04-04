package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheWebTek/four-market/utils/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Server interface {
	Start() error
}

type echoServer struct {
	echo   *echo.Echo
	config Config
	log    *logger.ZapSugaredLogger
}

func New(opts ...Option) Server {
	_ = godotenv.Load()

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Port == "" {
		cfg.Port = os.Getenv("PORT")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())

	return &echoServer{
		echo:   e,
		config: cfg,
		log:    logger.GetInstance(),
	}
}

func (s *echoServer) Start() error {
	s.log.Infow("starting server", "port", s.config.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startConfig := echo.StartConfig{
		Address:         ":" + s.config.Port,
		GracefulTimeout: s.config.GracefulTimeout,
		HideBanner:      true,
		HidePort:        true,
	}

	if err := startConfig.Start(ctx, s.echo); err != nil {
		s.log.Infow("server stopped", "error", err)
		return err
	}

	s.log.Info("server stopped")
	return nil
}

func Start() {
	New().Start()
}
