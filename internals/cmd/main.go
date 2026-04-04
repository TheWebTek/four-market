package main

import (
	"github.com/AbiXnash/event-manager/internals/server"
	logger "github.com/AbiXnash/event-manager/utils"
)

var log = logger.GetInstance()

func main() {
	server.Start()
	log.Info("Server Started!")
}
