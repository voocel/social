package group

import (
	"github.com/spf13/viper"
	"social/internal/node"
	"social/pkg/discovery"
)

func Run() *node.Node {
	n := node.NewNode(&discovery.Node{
		Name: "group",
		Host: viper.GetString("group.host"),
		Port: viper.GetInt("group.port"),
	})
	core := newCore(n.GetProxy())
	core.Init()
	n.Start()
	return n
}
