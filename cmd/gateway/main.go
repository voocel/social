package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"social/config"
	"social/internal/app/gateway"
	"social/internal/transport"
	"social/internal/transport/grpc"
	"social/pkg/log"
	"social/pkg/network/ws"
)

func main() {
	config.LoadConfig()
	log.Init("gateway", "debug")
	srv := ws.NewServer(viper.GetString("gateway.addr"))
	addr := viper.GetString("transport.grpc.addr")
	name := viper.GetString("transport.grpc.service_name")
	g := gateway.NewGateway(
		gateway.WithServer(srv),
		gateway.WithTransporter(grpc.NewTransporter(transport.WithAddr(addr), transport.WithName(name))),
	)
	g.Start()

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
