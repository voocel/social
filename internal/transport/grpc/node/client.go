package node

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"social/config"
	"social/internal/event"
	"social/internal/transport"
	"social/pkg/discovery/etcd"
	"social/pkg/log"
	"social/protos/pb"
)

var clients sync.Map

type client struct {
	client pb.NodeClient
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

	log.Infof("[Gateway] grpc client trying to connect to node[%s]...", serviceName)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", reg.Scheme(), serviceName), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithBlock())
	if err != nil {
		log.Errorf("[Gateway] the node[%s] grpc server not ready yet: %v", serviceName, err)
		return nil, err
	}

	log.Infof("[Gateway] grpc client connect to node[%s] is successful!", serviceName)

	cc := &client{client: pb.NewNodeClient(conn)}
	clients.Store(serviceName, cc)
	return cc, nil
}

func (c client) Trigger(ctx context.Context, event event.Event, gid string, uid int64) (err error) {
	_, err = c.client.Trigger(ctx, &pb.TriggerReq{
		Event: int32(event),
		Gid:   gid,
		Uid:   uid,
	})
	return err
}

func (c client) Deliver(ctx context.Context, cid, uid int64, message *transport.Message) (err error) {
	_, err = c.client.Deliver(ctx, &pb.DeliverReq{
		Cid: cid,
		Uid: uid,
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	})
	return err
}
