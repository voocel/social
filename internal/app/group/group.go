package group

import (
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"social/config"
	"social/internal/node"
	"social/internal/transport"
	"social/internal/transport/grpc"
	"social/pkg/database/ent"
)

func Run() *node.Node {
	addr := viper.GetString("transport.grpc.addr")
	name := viper.GetString("transport.grpc.service_name")
	n := node.NewNode(node.WithTransporter(
		grpc.NewTransporter(transport.WithAddr(addr), transport.WithName(name)),
	))
	cfg := ent.PgConfig{}
	copier.Copy(&cfg, config.Conf.Postgres)
	client := ent.InitEnt(cfg)
	core := newCore(n.GetProxy(), client)
	core.Init()
	n.Start()
	return n
}
