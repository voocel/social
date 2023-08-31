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

	client, ok := p.gate.nodeClient[serviceName]
	if !ok {
		return fmt.Errorf("node service client not found: %v", serviceName)
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
