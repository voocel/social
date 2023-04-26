package main

import (
	"social/config"
	"social/internal/app/im"
	"social/pkg/log"
)

func main() {
	config.LoadConfig()
	log.Init("im", "debug")
	im.Run()
}
