package discovery

import (
	"context"
	"fmt"
)

type Discovery interface {
	Name() string

	// Register 注册服务
	Register(context.Context, *Node, int64) error

	// Unregister 取消服务
	Unregister(context.Context, *Node) error

	// Query 查询指定服务
	Query(string) *Node

	// QueryServices 向注册中心查询所有服务
	QueryServices() []*Node
}

type Node struct {
	Id       string
	Name     string
	Host     string
	Port     int
	Enable   bool
	Healthy  bool
	Weight   float64
	Tags     []string
	Metadata map[string]string
}

func (n *Node) GetId() string {
	return n.Id
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) GetHost() string {
	return n.Host
}

func (n *Node) GetPort() int {
	return n.Port
}

func (n *Node) GetAddress() string {
	return fmt.Sprintf("%s:%d", n.GetHost(), n.GetPort())
}

func (n *Node) IsEnable() bool {
	return n.Enable
}

func (n *Node) IsHealthy() bool {
	return n.Healthy
}

func (n *Node) GetWeight() float64 {
	return n.Weight
}

func (n *Node) GetTags() []string {
	return n.Tags
}

func (n *Node) GetMetadata() map[string]string {
	if n.Metadata == nil {
		n.Metadata = make(map[string]string, 0)
	}
	return n.Metadata
}
