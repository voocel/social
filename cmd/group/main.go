package main

import (
	"social/config"
	"social/internal/app/group"
	"social/pkg/log"
)

func main() {
	config.LoadConfig()
	log.Init("group", "debug")
	group.Run()
}
