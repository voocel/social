package discovery

import "context"

type Discovery interface {
	Name() string

	// Register 注册服务
	Register(context.Context, Node, int64) error

	// UnRegister 取消服务
	UnRegister(context.Context, Node) error

	// QueryServices 向注册中心查询所有服务
	QueryServices() []*Node
}

type Node struct {
	Name string
	Host string
	Port string
	Meta map[string]string
}
