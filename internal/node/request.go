package node

import "context"

type Request struct {
	Cid    int64
	Uid    int64
	Seq    int32
	Route  int32
	Buffer []byte
	Node   *Node
}

func (r *Request) Respond(ctx context.Context, target int64, message []byte) error {
	return r.Node.proxy.Respond(ctx, r, target, message)
}
