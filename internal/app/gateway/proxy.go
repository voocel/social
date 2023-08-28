package gateway

import (
	"context"
	"social/protos/pb"
)

type proxy struct {
	gate *Gateway
	bind map[int64]string
}

func newProxy(gate *Gateway) *proxy {
	return &proxy{
		gate: gate,
	}
}

// Launch send to node
func (p *proxy) push(ctx context.Context, cid, uid int64, message []byte, route int32) ([]byte, error) {
	reply, err := p.gate.nodeClient["im"].Deliver(ctx, &pb.DeliverReq{
		Gid: p.gate.opts.id,
		Cid: cid,
		Uid: uid,
		Message: &pb.Message{
			Seq:    0,
			Route:  route,
			Buffer: message,
		},
	})
	return reply.GetPayload(), err
}

func (p *proxy) bindGate(ctx context.Context, uid int64) {
	p.bind[uid] = p.gate.opts.id
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate(ctx context.Context, uid int64) {
	delete(p.bind, uid)
}
