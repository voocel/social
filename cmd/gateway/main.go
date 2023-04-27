package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"social/config"
	"social/internal/app/gateway"
	"social/pkg/log"
	"social/pkg/network/ws"
	"syscall"
)

func main() {
	config.LoadConfig()
	log.Init("gateway", "debug")
	srv := ws.NewServer(fmt.Sprintf(viper.GetString("gateway.host")))
	g := gateway.NewGateway(gateway.WithServer(srv))
	g.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Sync()
			g.Stop()
		case syscall.SIGHUP:
			config.LoadConfig()
		default:
			return
		}
	}
}
