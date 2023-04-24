package main

import (
	"os"
	"os/signal"
	"social/config"
	"social/internal/app/gateway"
	"social/pkg/log"
	"syscall"
)

func main() {
	config.LoadConfig()
	log.Init("gateway", "log", "debug")
	g := gateway.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Sync()
			g.Stop()
			return
		case syscall.SIGHUP:
			config.LoadConfig()
		default:
			return
		}
	}
}
