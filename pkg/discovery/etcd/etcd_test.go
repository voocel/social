package etcd

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"social/pkg/discovery"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	d := NewDiscovery([]string{"127.0.0.1:2379"})
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
	d := NewDiscovery([]string{"127.0.0.1:2379"})
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

func TestGrpc(t *testing.T) {
	r := NewResolver([]string{"127.0.0.1:2379"}, "")
	resolver.Register(r)
	conn, err := grpc.DialContext(context.Background(), r.Scheme()+"://web", grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	fmt.Println(conn)
}
