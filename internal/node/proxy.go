package node

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"social/config"
	"social/internal/event"
	"social/internal/transport"
	"social/pkg/log"
)

type Proxy struct {
	node *Node
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

func (p *Proxy) getGateClient(name string) (transport.GateClient, error) {
	if len(name) == 0 {
		return nil, errors.New("gateway service name not be empty")
	}
	return p.node.opts.transporter.NewGateClient(name)
}

// Respond send to gateway grpc server
func (p *Proxy) Respond(ctx context.Context, req *Request, target int64, msg []byte) error {
	c, err := p.getGateClient(config.Conf.Transport.DiscoveryGate)
	if err != nil {
		return err
	}
	err = c.Push(ctx, target, &transport.Message{
		Seq:    0,
		Route:  0,
		Buffer: msg,
	})
	if err != nil {
		st := status.Convert(err)
		for _, d := range st.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				// User offline
				goto ok
			default:
				log.Errorf("Unexpected epb type: %v", info)
			}
		}
		return fmt.Errorf("send to gateway err: %v", err)
	}
ok:
	log.Infof("Respond message to gateway success: %s", string(msg))
	return nil
}
