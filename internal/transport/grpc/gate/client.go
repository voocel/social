package gate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"social/internal/transport"
	"social/pkg/discovery/etcd"
	"social/pkg/log"
	"social/protos/pb"
)

var clients sync.Map

type client struct {
	client pb.GateClient
}

func NewClient(serviceName string) (*client, error) {
	c, ok := clients.Load(serviceName)
	if ok {
		return c.(*client), nil
	}

	reg, err := etcd.NewResolver([]string{viper.GetString("etcd.addr")}, serviceName)
	if err != nil {
		panic(err)
	}
	resolver.Register(reg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infof("[Gateway] grpc client trying to connect to node [%s]...", serviceName)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", reg.Scheme(), serviceName), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithBlock())
	if err != nil {
		log.Warnf("[Gateway] the node[%s] grpc server not ready yet: %v", serviceName, err)
		return nil, err
	}

	cc := &client{client: pb.NewGateClient(conn)}
	clients.Store(serviceName, cc)
	return cc, nil
}

func (c *client) Bind(ctx context.Context, cid, uid int64) (err error) {
	_, err = c.client.Bind(ctx, &pb.BindReq{
		Cid: cid,
		Uid: uid,
	})

	return err
}

func (c *client) Unbind(ctx context.Context, uid int64) (err error) {
	_, err = c.client.Unbind(ctx, &pb.UnbindReq{
		Uid: uid,
	})
	return err
}

func (c *client) GetIP(ctx context.Context, target int64) (ip string, err error) {
	res, e := c.client.GetIP(ctx, &pb.GetIPReq{
		Uid: target,
	})
	if e != nil {
		return "", e
	}
	return res.GetIP(), nil
}

func (c *client) Disconnect(ctx context.Context, target int64) (err error) {
	panic("implement me")
}

// Push node send message to gateway grpc server
func (c *client) Push(ctx context.Context, target int64, message *transport.Message) (err error) {
	_, err = c.client.Push(ctx, &pb.PushReq{
		Target: target,
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	})
	return err
}

func (c *client) Broadcast(ctx context.Context, message *transport.Message) (total int64, err error) {
	panic("implement me")
}
