package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"social/pkg/discovery"
	"social/pkg/discovery/etcd"
	"social/pkg/log"
	"social/protos/node"
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
	handler  RouteHandler
}

type Node struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	proxy               *Proxy
	opts                *options
	registry            *etcd.Registry
	eventCh             chan *eventEntity
	routes              map[int32]routeEntity
	events              map[Event]EventHandler
	instance            *discovery.Node
	rpcSrv              *grpc.Server
	defaultRouteHandler RouteHandler
	sync.RWMutex
}

func NewNode(instance *discovery.Node, opts ...OptionFunc) *Node {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	n := &Node{}
	n.proxy = newProxy(n)
	n.instance = instance
	n.routes = make(map[int32]routeEntity)
	n.events = make(map[Event]EventHandler)
	n.eventCh = make(chan *eventEntity, 1024)
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
			handler, ok := n.events[entity.event]
			if !ok {
				log.Warnf("event does not register handler function, event: %v", entity.event)
				continue
			}
			handler(entity.gid, entity.uid)
		}
	}
}

func (n *Node) Stop() {
	if err := n.registry.Unregister(n.ctx, n.instance); err != nil {
		log.Errorf("[%s]registry unregister err: %v", n.instance.Name, err)
	}
	n.rpcSrv.GracefulStop()
	close(n.eventCh)
	n.cancel()
	log.Infof("[node] stop and unregister successful: %v", n.instance.Name)
}

func (n *Node) startRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.instance.Host, n.instance.Port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	n.rpcSrv = s
	node.RegisterNodeServer(s, &nodeService{node: n})

	r, err := etcd.NewRegistry([]string{viper.GetString("etcd.addr")})
	if err != nil {
		panic(err)
	}
	n.registry = r

	err = n.registry.Register(n.ctx, n.instance, 60)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	log.Infof("[%s]node rpc server stop successful", n.instance.Name)
}

func (n *Node) addRouteHandler(route int32, handler RouteHandler) {
	n.Lock()
	defer n.Unlock()
	n.routes[route] = routeEntity{
		route:    route,
		stateful: false,
		handler:  handler,
	}
}

func (n *Node) addEventListener(event Event, handler EventHandler) {
	n.events[event] = handler
}

func (n *Node) triggerEvent(event Event, gid string, uid int64) {
	n.eventCh <- &eventEntity{
		event: event,
		gid:   gid,
		uid:   uid,
	}
}
