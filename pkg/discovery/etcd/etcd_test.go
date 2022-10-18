package etcd

import (
	"context"
	"fmt"
	"social/pkg/discovery"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	d := NewDiscovery()
	defer d.Close()

	d.Watch("/grpc/")
	d.Watch("/web/")
	for {
		select {
		case <-time.Tick(5 * time.Second):
			res := d.QueryServices()
			for _, v := range res {
				fmt.Println(v)
			}
		}
	}
}

func TestRegistry(t *testing.T) {
	d := NewDiscovery()
	d.Registry.Register(context.Background(), discovery.Node{
		Name: "web",
		Host: "node1",
		Port: "8888",
	}, 10)
	d.Registry.Register(context.Background(), discovery.Node{
		Name: "web",
		Host: "node2",
		Port: "8888",
	}, 10)
	d.Registry.Register(context.Background(), discovery.Node{
		Name: "grpc",
		Host: "node1",
		Port: "8888",
	}, 10)

	select {
	case <-time.After(60 * time.Second):
		d.Registry.Close()
	}
}
