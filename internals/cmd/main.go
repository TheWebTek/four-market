package main

import (
	"github.com/TheWebTek/four-market/internals/server"
	"github.com/TheWebTek/four-market/utils/logger"
)

func main() {
	log := logger.GetInstance()
	log.Info("Server Starting...")
	server.Start()
}
