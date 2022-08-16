package main

import (
	"social/internal/app/gateway"
	"social/pkg/log"
)

func main()  {
	log.Init("gateway", "log", "debug")
	gateway.Run()
}