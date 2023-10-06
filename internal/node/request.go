package node

import (
	"context"
	"social/protos/pb"
)

type Request struct {
	Cid    int64
	Uid    int64
	Seq    int32
	Route  int32
	Buffer []byte
	Node   *Node
}

func (r *Request) Respond(ctx context.Context, target int64, message *pb.MsgItem) error {
	return r.Node.proxy.Respond(ctx, r, target, message)
}

func (r *Request) Multicast(ctx context.Context, target int64, message *pb.MsgItem) error {
	return r.Node.proxy.Multicast(ctx, target, message)
}

func (r *Request) Broadcast(ctx context.Context, target int64, message *pb.MsgItem) error {
	return r.Node.proxy.Broadcast(ctx, message)
}
