package node

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"social/internal/event"
	"social/internal/transport"
	"social/protos/pb"
	"sync"
)

var clients sync.Map

type client struct {
	client pb.NodeClient
}

func NewClient(addr string) (*client, error) {
	c, ok := clients.Load(addr)
	if ok {
		return c.(*client), nil
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cc := &client{client: pb.NewNodeClient(conn)}
	clients.Store(addr, cc)
	return cc, nil
}

func (c client) Trigger(ctx context.Context, event event.Event, gid string, uid int64) (err error) {
	_, err = c.client.Trigger(ctx, &pb.TriggerRequest{
		Event: int32(event),
		Gid:   gid,
		Uid:   uid,
	})
	return err
}

func (c client) Deliver(ctx context.Context, gid, nid string, cid, uid int64, message *transport.Message) (err error) {
	_, err = c.client.Deliver(ctx, &pb.DeliverRequest{
		Gid: gid,
		Nid: nid,
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
