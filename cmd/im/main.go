package main

import (
	"os"
	"os/signal"
	"social/config"
	"social/internal/app/im"
	"social/pkg/log"
	"social/pkg/redis"
	"syscall"
)

func main() {
	config.LoadConfig()
	log.Init("im", "debug")
	n := im.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			n.Stop()
			log.Sync()
			redis.Close()
			return
		case syscall.SIGHUP:
			config.LoadConfig()
		default:
			return
		}
	}
}
