package gateway

import (
	"context"
	"social/protos/node"
)

type proxy struct {
	gate       *Gateway
	nodeClient node.NodeClient
}

func newProxy(gate *Gateway) *proxy {
	return &proxy{
		gate: gate,
	}
}

// Launch send to node
func (p *proxy) push(ctx context.Context, cid, uid int64, message []byte, route int32) ([]byte, error) {
	reply, err := p.nodeClient.Deliver(ctx, &node.DeliverRequest{
		Gid:    p.gate.opts.id,
		Cid:    cid,
		Uid:    uid,
		Route:  route,
		Buffer: message,
	})
	return reply.GetPayload(), err
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate() {

}

func (p *proxy) newNodeClient(name string) {
	p.nodeClient = node.NewNodeClient(p.gate.nodeConns[name])
}
