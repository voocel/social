package gate

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"social/internal/transport"
	"social/protos/pb"
	"sync"
)

var clients sync.Map

type client struct {
	client pb.GateClient
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
	cc := &client{client: pb.NewGateClient(conn)}
	clients.Store(addr, cc)
	return cc, nil
}

func (c *client) Bind(ctx context.Context, cid, uid int64) (err error) {
	_, err = c.client.Bind(ctx, &pb.BindRequest{
		Cid: cid,
		Uid: uid,
	})

	return err
}

func (c *client) Unbind(ctx context.Context, uid int64) (err error) {
	_, err = c.client.Unbind(ctx, &pb.UnbindRequest{
		Uid: uid,
	})
	return err
}

func (c *client) GetIP(ctx context.Context, target int64) (ip string, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) Disconnect(ctx context.Context, target int64) (err error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) Push(ctx context.Context, target int64, message *transport.Message) (err error) {
	_, err = c.client.Push(ctx, &pb.PushRequest{
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
	//TODO implement me
	panic("implement me")
}

func (c *client) Broadcast(ctx context.Context, message *transport.Message) (total int64, err error) {
	//TODO implement me
	panic("implement me")
}
