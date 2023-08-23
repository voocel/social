package node

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"social/internal/event"
	"social/internal/router"
	"social/internal/transport"
	"social/pkg/discovery"
	"social/pkg/discovery/etcd"
	"social/pkg/log"
)

type RouteHandler func(req Request) error

type Request struct {
	Gid    string
	Nid    string
	Cid    int64
	Uid    int64
	Route  int32
	Buffer []byte
	Node   *Node
}

type routeEntity struct {
	route    int32
	stateful bool
	Handler  RouteHandler
}

type Node struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	proxy               *Proxy
	opts                *options
	registry            *etcd.Registry
	eventCh             chan *event.EventEntity
	router              *router.Router
	Routes              map[int32]routeEntity
	events              map[event.Event]event.EventHandler
	instance            *discovery.Node
	rpcSrv              transport.Server
	transporter         transport.Transporter
	DefaultRouteHandler RouteHandler
	sync.RWMutex
}

func NewNode(opts ...OptionFunc) *Node {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	n := &Node{}
	n.opts = o
	n.proxy = newProxy(n)
	n.router = router.NewRouter()
	n.Routes = make(map[int32]routeEntity)
	n.events = make(map[event.Event]event.EventHandler)
	n.eventCh = make(chan *event.EventEntity, 1024)
	n.ctx, n.cancel = context.WithCancel(context.Background())
	return n
}

func (n *Node) GetProxy() *Proxy {
	return n.proxy
}

func (n *Node) Start() {
	go n.startRPCServer()
	go n.dispatch()
}

func (n *Node) dispatch() {
	for {
		select {
		case entity, ok := <-n.eventCh:
			if !ok {
				return
			}
			handler, ok := n.events[entity.Event]
			if !ok {
				log.Warnf("event does not register handler function, event: %v", entity.Event)
				continue
			}
			handler(entity.Gid, entity.Uid)
		}
	}
}

func (n *Node) Stop() {
	if err := n.registry.Unregister(n.ctx, n.instance); err != nil {
		log.Errorf("[%s]registry unregister err: %v", n.instance.Name, err)
	}
	n.rpcSrv.Stop()
	close(n.eventCh)
	n.cancel()
	log.Infof("[node] stop and unregister successful: %v", n.instance.Name)
}

func (n *Node) startRPCServer() {
	n.rpcSrv = n.opts.transporter.NewNodeServer(&provider{n})

	r, err := etcd.NewRegistry([]string{viper.GetString("etcd.addr")})
	if err != nil {
		panic(err)
	}
	n.registry = r

	instance := &discovery.Node{
		Name: n.opts.transporter.Options().Name,
		Addr: n.opts.transporter.Options().Server.Addr,
	}
	n.instance = instance

	err = n.registry.Register(n.ctx, instance, 60)
	if err != nil {
		log.Fatal(err)
	}

	if err := n.rpcSrv.Start(); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	log.Infof("[%s]node rpc server stop successful", n.instance.Name)
}

func (n *Node) addRouteHandler(route int32, handler RouteHandler) {
	n.Lock()
	defer n.Unlock()
	n.Routes[route] = routeEntity{
		route:    route,
		stateful: false,
		Handler:  handler,
	}
}

func (n *Node) addEventListener(event event.Event, handler event.EventHandler) {
	n.events[event] = handler
}

func (n *Node) triggerEvent(e event.Event, gid string, uid int64) {
	n.eventCh <- &event.EventEntity{
		Event: e,
		Gid:   gid,
		Uid:   uid,
	}
}
