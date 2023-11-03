package im

import (
	"github.com/jinzhu/copier"
	"social/config"
	"social/internal/node"
	"social/internal/transport"
	"social/internal/transport/grpc"
	"social/pkg/database/ent"
	"social/pkg/redis"
)

func Run() *node.Node {
	addr := config.Conf.Transport.Grpc.Addr
	name := config.Conf.Transport.Grpc.ServiceName
	n := node.NewNode(
		node.WithTransporter(
			grpc.NewTransporter(transport.WithAddr(addr), transport.WithName(name)),
		),
	)
	cfg := ent.PgConfig{}
	copier.Copy(&cfg, config.Conf.Postgres)
	client := ent.InitEnt(cfg)
	core := newCore(n.GetProxy(), client)
	core.Init()
	redis.Init()
	n.Start()
	return n
}
