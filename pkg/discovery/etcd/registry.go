package etcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Registry struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	key           string
	val           string
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func NewRegistry(endpoints []string, key, val string, lease int64) (*Registry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	r := &Registry{
		cli:           cli,
		key:           key,
		val:           val,
	}
	if err = r.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) putKeyWithLease(lease int64) error {
	//设置租约时间
	resp, err := r.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = r.cli.Put(context.Background(), r.key, r.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租定期发送需求请求
	leaseResp, err := r.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	r.leaseID = resp.ID
	r.keepAliveChan = leaseResp
	log.Printf("put key:%s  value:%s  success!", r.key, r.val)
	return nil
}

//ListenLeaseChan 监听续租情况
func (r *Registry) ListenLeaseChan() {
	for leaseKeepResp := range r.keepAliveChan {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("关闭续租")
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