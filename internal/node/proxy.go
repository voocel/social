package node

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"social/internal/transport"
	"social/protos/pb"
	"sync"
)

var gateClients sync.Map

type gateClient struct {
	client pb.GateClient
}

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

func (p *Proxy) AddEventListener(event Event, handler EventHandler) {
	p.node.addEventListener(event, handler)
}

func (p *Proxy) BindGate(gid string, cid, uid int64) error {
	c := p.getGateClient("127.0.0.1:7400")
	_, err := c.client.Bind(context.Background(), &pb.BindRequest{
		Cid: cid,
		Uid: uid,
	})
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

func (p *Proxy) getGateClient(addr string) *gateClient {
	c, ok := gateClients.Load(addr)
	if ok {
		return c.(*gateClient)
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	cc := &gateClient{client: pb.NewGateClient(conn)}
	gateClients.Store(addr, cc)
	return cc
}

func (p *Proxy) getGateClientByGid(gid string) (transport.GateClient, error) {
	endpoint, err := p.node.router.FindGateEndpoint(gid)
	if err != nil {
		return nil, err
	}
	return p.node.transporter.NewGateClient(endpoint)
}

func (p *Proxy) Respond(ctx context.Context, target int64, msg []byte) error {
	c, err := p.getGateClientByGid("")
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
