package main

import (
	"social/config"
	"social/internal/app/gateway"
	"social/pkg/log"
)

func main() {
	config.LoadConfig()
	log.Init("gateway", "log", "debug")
	gateway.Run()
}
