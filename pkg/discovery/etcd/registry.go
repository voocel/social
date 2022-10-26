package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"social/pkg/discovery"
)

const defaultRefreshDuration = time.Second * 10

type Registry struct {
	cc              resolver.ClientConn
	cli             *clientv3.Client
	leaseID         clientv3.LeaseID
	refreshDuration time.Duration
	serviceList     map[string]*discovery.Node
	sync.Mutex
}

func NewRegistry(endpoints []string) (*Registry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	r := &Registry{
		cli:         cli,
		serviceList: make(map[string]*discovery.Node),
	}
	return r, nil
}

func (r *Registry) Name() string {
	return "etcd"
}

func (r *Registry) Register(ctx context.Context, info *discovery.Node, lease int64) error {
	key := r.key(info)
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	value := string(b)
	//设置租约时间
	leaseResp, err := r.cli.Grant(ctx, lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = r.cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	//定期刷新租约使其不过期
	leaseKeepResp, err := r.cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}
	r.leaseID = leaseResp.ID

	// 监听续租情况
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Println("关闭续租")
				return
			case _, ok := <-leaseKeepResp:
				if !ok {
					log.Println("关闭续租")
					return
				}
				log.Println("续约成功", leaseKeepResp)
			}
		}
	}(ctx)
	log.Printf("put key:%s  value:%s  success!", key, value)
	return nil
}

func (r *Registry) Unregister(ctx context.Context, info *discovery.Node) error {
	_, err := r.cli.Delete(ctx, r.key(info))
	return err
}

func (r *Registry) QueryServices() []*discovery.Node {
	addrs := make([]*discovery.Node, 0, 10)
	for _, node := range r.serviceList {
		addrs = append(addrs, node)
	}
	return addrs
}

func (r *Registry) watch(keyPrefix string) error {
	if keyPrefix == "" {
		return errors.New("serviceName is empty")
	}

	r.setServices(keyPrefix)
	go r.watcher(keyPrefix)
	go r.refresh(keyPrefix)
	return nil
}

func (r *Registry) watcher(prefix string) {
	addrMap := make(map[string]resolver.Address)
	updateStateFunc := func() {
		addrList := make([]resolver.Address, 0, len(addrMap))
		for _, v := range addrMap {
			addrList = append(addrList, v)
		}
		log.Printf("watch addr: %v", addrList)
		r.cc.UpdateState(resolver.State{Addresses: addrList})
	}
	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i := range resp.Kvs {
			addr := strings.TrimPrefix(string(resp.Kvs[i].Key), prefix)
			addrMap[addr] = resolver.Address{Addr: addr}
		}
	}
	updateStateFunc()
	watchCh := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for ch := range watchCh {
		for _, event := range ch.Events {
			key := string(event.Kv.Key)
			addr := strings.TrimPrefix(key, prefix)
			switch event.Type {
			case mvccpb.PUT:
				info := &discovery.Node{}
				json.Unmarshal(event.Kv.Value, info)
				addrMap[addr] = resolver.Address{Addr: addr}
				r.setService(key, info)
				log.Println("PUT: ", key)
			case mvccpb.DELETE:
				delete(addrMap, addr)
				r.delService(key)
				log.Println("DELETE: ", key)
			}
		}
		updateStateFunc()
	}
}

func (r *Registry) refresh(prefix string) {
	if r.refreshDuration == -1 {
		return
	}
	if r.refreshDuration == 0 {
		r.refreshDuration = defaultRefreshDuration
	}
	ticker := time.NewTicker(r.refreshDuration)
	for range ticker.C {
		r.setServices(prefix)
		log.Println("refresh all")
	}
}

func (r *Registry) setServices(prefix string) {
	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Printf("get by prefix [%v] err: %v", prefix, err)
		return
	}
	for _, kv := range resp.Kvs {
		info := &discovery.Node{}
		json.Unmarshal(kv.Value, info)
		r.setService(string(kv.Key), info)
	}
}

func (r *Registry) setService(k string, v *discovery.Node) {
	r.Lock()
	defer r.Unlock()
	r.serviceList[k] = v
	log.Println("put key :", k, "value:", v)
}

func (r *Registry) delService(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.serviceList, key)
}

func (r *Registry) Close() error {
	//撤销租约
	if _, err := r.cli.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return r.cli.Close()
}

func (r *Registry) key(info *discovery.Node) string {
	return fmt.Sprintf("etcd/%s/%s", info.GetName(), info.GetAddress())
}
