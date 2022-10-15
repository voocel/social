package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"social/pkg/discovery"
	"time"
)

type Registry struct {
	cli     *clientv3.Client
	leaseID clientv3.LeaseID
}

func NewRegistry() (*Registry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	r := &Registry{
		cli: cli,
	}
	return r, nil
}

func (r *Registry) Register(ctx context.Context, info discovery.Node, lease int64) error {
	key := fmt.Sprintf("/%s/%s", info.Name, info.ID)
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

func (r *Registry) UnRegister(ctx context.Context, info discovery.Node) error {
	key := fmt.Sprintf("/%s/%s", info.Name, info.ID)
	_, err := r.cli.Delete(ctx, key)
	return err
}

// Close 注销服务
func (r *Registry) Close() error {
	//撤销租约
	if _, err := r.cli.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return r.cli.Close()
}
