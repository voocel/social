package im

import (
	"github.com/spf13/viper"
	"social/internal/node"
	"social/pkg/discovery"
)

func Run() *node.Node {
	n := node.NewNode(&discovery.Node{
		Name: "im",
		Host: viper.GetString("im.host"),
		Port: viper.GetInt("im.port"),
	})
	core := newCore(n.GetProxy())
	core.Init()
	n.Start()
	return n
}
