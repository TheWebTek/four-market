package server

import (
	"context"
	"os"
	"strings"

	logger "github.com/AbiXnash/event-manager/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var log = logger.GetInstance()

func Start() {
	port := getServerPort()
	spawnServer(port)
}

func getServerPort() string {
	godotenv.Load()
	port := strings.TrimSpace(os.Getenv("PORT"))

	log.Infow("Starting server in ", "port", port)

	return ":" + port
}

func spawnServer(port string) {
	e := echo.New()
	e.Use(middleware.RequestLogger())

	startConfig := echo.StartConfig{
		Address:    port,
		HideBanner: true,
		HidePort:   true,
	}

	if err := startConfig.Start(context.Background(), e); err != nil {
		log.Info("failed to start server", "error", err)
	}
}
