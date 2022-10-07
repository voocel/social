package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 10*time.Second,
	})
	if err != nil {
		panic(err)
	}
	cli.Put(context.Background(), "test", "hello")
	resp, err := cli.Get(context.Background(), "test")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Kvs)
	res, err := cli.Grant(context.Background(), 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ID)
	cli.Put(context.Background(), "test2", "666", clientv3.WithLease(res.ID))
	for i := 0; i < 15; i++ {
		r , _ := cli.Get(context.Background(), "test2")
		fmt.Println(r)
		time.Sleep(time.Second)
	}
}

func TestDiscovery(t *testing.T) {
	d := NewDiscovery()
	defer d.Close()
	d.Watch("/web/")
	d.Watch("/grpc/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(d.GetServices())
		}
	}
}

func TestRegistry(t *testing.T) {
	var endpoints = []string{"node1:6666"}
	ser, err := NewRegistry(endpoints, "/web/node1", "localhost:2379", 20)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	go ser.ListenLeaseChan()
	select {
	case <-time.After(30 * time.Second):
		ser.Close()
	}
}
