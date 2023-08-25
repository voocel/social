package im

import (
	"github.com/spf13/viper"
	"social/internal/node"
	"social/internal/transport"
	"social/internal/transport/grpc"
)

func Run() *node.Node {
	addr := viper.GetString("transport.grpc.addr")
	name := viper.GetString("transport.grpc.service_name")
	n := node.NewNode(
		node.WithTransporter(
			grpc.NewTransporter(transport.WithAddr(addr), transport.WithName(name)),
		),
	)
	core := newCore(n.GetProxy())
	core.Init()
	n.Start()
	return n
}
