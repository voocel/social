package main

import (
	"social/config"
	"social/internal/app/http"
	"social/pkg/log"
)

func main() {
	config.LoadConfig()
	log.Init("http", "debug")
	srv := http.NewServer()
	srv.Run()
}
