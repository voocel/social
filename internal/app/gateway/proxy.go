package gateway

import (
	"context"
	"social/config"
	"social/internal/transport"
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
func (p *proxy) push(ctx context.Context, cid, uid int64, message []byte, route int32) error {
	var serviceName string
	for _, v := range config.Conf.Transport.Discovery {
		if inSlice(v.Routers, route) {
			serviceName = v.Name
		}
	}

	return p.gate.nodeClient[serviceName].Deliver(ctx, cid, uid,
		&transport.Message{
			Seq:    0,
			Route:  route,
			Buffer: message,
		},
	)
}

func (p *proxy) bindGate(ctx context.Context, uid int64) {
	p.bind[uid] = p.gate.opts.id
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate(ctx context.Context, uid int64) {
	delete(p.bind, uid)
}

func inSlice(s []int32, v int32) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}
