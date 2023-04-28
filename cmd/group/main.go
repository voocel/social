package main

import (
	"os"
	"os/signal"
	"social/config"
	"social/internal/app/group"
	"social/pkg/log"
	"syscall"
)

func main() {
	config.LoadConfig()
	log.Init("group", "debug")
	n := group.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			n.Stop()
			log.Sync()
		case syscall.SIGHUP:
			config.LoadConfig()
		default:
			return
		}
	}
}
