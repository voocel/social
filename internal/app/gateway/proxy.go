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

// Launch
func (p *proxy) push(ctx context.Context, cid, uid int64, message []byte, route int32) error {
	_, err := p.nodeClient.Deliver(ctx, &node.DeliverRequest{
		GID:    "",
		CID:    cid,
		UID:    uid,
		Route:  route,
		Buffer: message,
	})
	return err
}

func (p *proxy) newNodeClient(name string) {
	p.nodeClient = node.NewNodeClient(p.gate.nodeConns[name])
}
