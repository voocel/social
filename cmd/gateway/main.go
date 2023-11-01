package main

import (
	"os"
	"os/signal"
	"syscall"

	"social/config"
	"social/internal/app/gateway"
	"social/internal/transport"
	"social/internal/transport/grpc"
	"social/pkg/log"
	"social/pkg/network/ws"
)

func main() {
	config.LoadConfig()
	log.Init(config.Conf.Name, config.Conf.LogLevel)

	srv := ws.NewServer(config.Conf.Gateway.Addr)
	addr := config.Conf.Transport.Grpc.Addr
	name := config.Conf.Transport.Grpc.ServiceName
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
