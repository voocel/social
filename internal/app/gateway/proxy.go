package gateway

import (
	"context"
	"fmt"
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

func (p *proxy) getNodeClient(name string) (transport.NodeClient, error) {
	return p.gate.opts.transporter.NewNodeClient(name)
}

// Launch send to node rpc
func (p *proxy) push(ctx context.Context, cid, uid int64, route int32, message []byte) error {
	var serviceName string
	for _, v := range config.Conf.Transport.DiscoveryNode {
		if inSlice(v.Routers, route) {
			serviceName = v.Name
		}
	}
	if len(serviceName) == 0 {
		return fmt.Errorf("service name not found: %v", serviceName)
	}

	client, err := p.getNodeClient(serviceName)
	if err != nil {
		return fmt.Errorf("[%v]gateway get node rpc client err: %v", serviceName, err)
	}

	return client.Deliver(ctx, cid, uid,
		&transport.Message{
			Seq:    0,
			Route:  route,
			Buffer: message,
		},
	)
}

func inSlice(s []int32, v int32) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}
