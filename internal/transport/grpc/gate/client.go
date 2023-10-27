package gate

import (
	"context"
	"fmt"
	"social/config"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
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

	reg, err := etcd.NewResolver([]string{config.Conf.Etcd.Addr}, serviceName)
	if err != nil {
		panic(err)
	}
	resolver.Register(reg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infof("[Node] grpc client trying to connect to node [%s]...", serviceName)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", reg.Scheme(), serviceName), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithBlock())
	if err != nil {
		log.Warnf("[Node] the gateway[%s] grpc server not ready yet: %v", serviceName, err)
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
	_, err = c.client.Unbind(ctx, &pb.BaseReq{
		Uid: uid,
	})
	return err
}

func (c *client) GetIP(ctx context.Context, target int64) (ip string, err error) {
	res, e := c.client.GetIP(ctx, &pb.BaseReq{
		Uid: target,
	})
	if e != nil {
		return "", e
	}
	return res.GetIp(), nil
}

func (c *client) Disconnect(ctx context.Context, target int64) (err error) {
	_, err = c.client.Disconnect(ctx, &pb.BaseReq{Uid: target})
	return err
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

func (c *client) Multicast(ctx context.Context, targets []int64, message *transport.Message) (total int64, err error) {
	reply, err := c.client.Multicast(ctx, &pb.MulticastReq{
		Targets: targets,
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	})
	return reply.Total, err
}

func (c *client) Broadcast(ctx context.Context, message *transport.Message) (total int64, err error) {
	reply, err := c.client.Broadcast(ctx, &pb.BroadcastReq{
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	}, grpc.UseCompressor(gzip.Name))
	return reply.Total, err
}
