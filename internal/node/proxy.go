package node

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"social/protos/gate"
	"sync"
)

var gateClients sync.Map

type gateClient struct {
	client gate.GateClient
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
	p.node.defaultRouteHandler = handler
}

func (p *Proxy) AddEventListener(event Event, handler EventHandler) {
	p.node.addEventListener(event, handler)
}

func (p *Proxy) BindGate(gid string, cid, uid int64) {
	c := p.getGateClient("")
	c.client.Bind(context.Background(), &gate.BindRequest{
		Cid: cid,
		Uid: uid,
	})
}

func (p *Proxy) getGateClient(addr string) *gateClient {
	c, ok := gateClients.Load(addr)
	if ok {
		return c.(*gateClient)
	}
	conn, err := grpc.Dial("127.0.0.1:7400", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	cc := &gateClient{client: gate.NewGateClient(conn)}
	gateClients.Store(addr, cc)
	return cc
}

func (p *Proxy) Respond(ctx context.Context, target int64, msg []byte) {
	c := p.getGateClient("")
	r, err := c.client.Push(ctx, &gate.PushRequest{Target: target, Buffer: msg})
	if err != nil {
		log.Fatalf("could not gate: %v", err)
	}
	log.Printf("push: %s", r.String())
}

func (p *Proxy) directDeliver(ctx context.Context, args *DeliverArgs) {

}
