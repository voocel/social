package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"log"
	"social/pkg/discovery"
	"sync"
	"time"
)

const defaultRefreshDuration = time.Second * 10

type Discovery struct {
	cli             *clientv3.Client
	serverList      map[string]*discovery.Node
	refreshDuration time.Duration
	serviceName     string
	Registry        *Registry
	sync.Mutex
}

func NewDiscovery() *Discovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	reg, err := NewRegistry()
	return &Discovery{
		cli:      cli,
		Registry: reg,
	}
}

func (d *Discovery) Register(ctx context.Context, info discovery.Node, lease int64) error {
	return d.Registry.Register(ctx, info, lease)
}

func (d *Discovery) UnRegister(ctx context.Context, info discovery.Node) error {
	return d.Registry.UnRegister(ctx, info)
}

func (d *Discovery) QueryServices() []*discovery.Node {
	addrs := make([]*discovery.Node, 0, 10)
	for _, node := range d.serverList {
		addrs = append(addrs, node)
	}
	return addrs
}

func (d *Discovery) Watch(keyPrefix string) error {
	if keyPrefix == "" {
		return errors.New("serviceName is empty")
	}
	d.serviceName = keyPrefix

	d.setServices()
	go d.watcher(keyPrefix)
	go d.refresh()
	return nil
}

func (d *Discovery) watcher(prefix string) {
	watchCh := d.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for ch := range watchCh {
		for _, event := range ch.Events {
			key := string(event.Kv.Key)
			switch event.Type {
			case mvccpb.PUT:
				info := &discovery.Node{}
				json.Unmarshal(event.Kv.Value, info)
				d.setService(key, info)
				log.Println("PUT: ", key)
			case mvccpb.DELETE:
				d.delService(key)
				log.Println("DELETE: ", key)
			}
		}
	}
}

func (d *Discovery) refresh() {
	if d.refreshDuration == -1 {
		return
	}
	if d.refreshDuration == 0 {
		d.refreshDuration = defaultRefreshDuration
	}
	ticker := time.NewTicker(d.refreshDuration)
	for range ticker.C {
		d.setServices()
		log.Println("refresh all!")
	}
}

func (d *Discovery) setServices() {
	resp, err := d.cli.Get(context.Background(), d.serviceName, clientv3.WithPrefix())
	if err != nil {
		log.Printf("get by prefix [%v] err: %v", d.serviceName, err)
		return
	}
	for _, kv := range resp.Kvs {
		info := &discovery.Node{}
		json.Unmarshal(kv.Value, info)
		d.setService(string(kv.Key), info)
	}
}

func (d *Discovery) setService(k string, v *discovery.Node) {
	d.Lock()
	d.Unlock()
	d.serverList[k] = v
	log.Println("put key :", k, "value:", v)
}

func (d *Discovery) delService(key string) {
	d.Lock()
	defer d.Unlock()
	delete(d.serverList, key)
}

func (d *Discovery) Close() error {
	return d.cli.Close()
}
