package node

import (
	"context"
	"fmt"
	"log"
	"social/internal/event"
	"social/internal/transport"
	"sync"
)

type Proxy struct {
	node       *Node
	sourceGate sync.Map
	sourceNode sync.Map
}

func newProxy(node *Node) *Proxy {
	return &Proxy{node: node}
}

func (p *Proxy) AddRouteHandler(route int32, handler RouteHandler) {
	p.node.addRouteHandler(route, handler)
}

// SetDefaultRouteHandler 设置默认路由处理器，所有未注册的路由均走默认路由处理器
func (p *Proxy) SetDefaultRouteHandler(handler RouteHandler) {
	p.node.DefaultRouteHandler = handler
}

func (p *Proxy) AddEventListener(event event.Event, handler event.EventHandler) {
	p.node.addEventListener(event, handler)
}

func (p *Proxy) BindGate(gid string, cid, uid int64) error {
	c, err := p.getGateClientByGid(gid)
	err = c.Bind(context.Background(), cid, uid)
	p.sourceGate.Store(uid, gid)
	return err
}

func (p *Proxy) GetGidByUid(uid int64) string {
	if val, ok := p.sourceGate.Load(uid); ok {
		if gid := val.(string); gid != "" {
			return gid
		}
	}
	return ""
}

func (p *Proxy) getGateClientByGid(gid string) (transport.GateClient, error) {
	return p.node.opts.transporter.NewGateClient(gid)
}

func (p *Proxy) getGateClientByUid(uid int64) (transport.GateClient, error) {
	if val, ok := p.sourceGate.Load(uid); ok {
		if gid := val.(string); gid != "" {
			return p.node.opts.transporter.NewGateClient(gid)
		}
	}
	return nil, fmt.Errorf("gateway not found by uid: %v", uid)
}

func (p *Proxy) Respond(ctx context.Context, req *Request, target int64, msg []byte) error {
	c, err := p.getGateClientByGid("gate-inter-rpc-server")
	if err != nil {
		return err
	}
	err = c.Push(ctx, target, &transport.Message{
		Seq:    0,
		Route:  0,
		Buffer: msg,
	})
	if err != nil {
		log.Fatalf("could not gate: %v", err)
	}
	log.Printf("push: %s", string(msg))
	return nil
}

func (p *Proxy) directDeliver(ctx context.Context, args *DeliverArgs) {

}
