package etcd

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

type Discovery struct {
	cli        *clientv3.Client
	serverList sync.Map
}

func NewDiscovery() *Discovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return &Discovery{
		cli: cli,
	}
}

func (d *Discovery) Watch(prefix string) error {
	resp, err := d.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		d.SetService(string(kv.Key), string(kv.Value))
	}
	go d.watcher(prefix)
	return nil
}

func (d *Discovery) watcher(prefix string) {
	watchCh := d.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for ch := range watchCh {
		for _, event := range ch.Events {
			switch event.Type {
			case mvccpb.DELETE:
				d.DelService(string(event.Kv.Key))
			case mvccpb.PUT:
				d.SetService(string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}
}

func (d *Discovery) SetService(k, v string) {
	d.serverList.Store(k, v)
	log.Println("put key :", k, "value:", v)
}

func (d *Discovery) DelService(key string) {
	d.serverList.Delete(key)
}

func (d *Discovery) GetServices() []string {
	addrs := make([]string, 0, 10)
	d.serverList.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(string))
		return true
	})
	return addrs
}

func (d *Discovery) Close() error {
	return d.cli.Close()
}
