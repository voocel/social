package node

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"social/pkg/discovery"
	"social/pkg/discovery/etcd"
	"social/pkg/log"
	"social/protos/node"
	"sync"
)

type RouteHandler func(req Request)

type Request struct {
	Gid    string
	Nid    string
	Cid    int64
	Uid    int64
	Route  int32
	Buffer interface{}
	Node   *Node
}

type routeEntity struct {
	route    int32
	stateful bool
	handler  RouteHandler
}

type Node struct {
	ctx                 context.Context
	proxy               *Proxy
	opts                *options
	registry            *etcd.Registry
	routes              map[int32]routeEntity
	defaultRouteHandler RouteHandler
	sync.RWMutex
}

func NewNode(opts ...OptionFunc) *Node {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	n := &Node{}
	n.proxy = newProxy(n)
	n.routes = make(map[int32]routeEntity)
	return n
}

func (n *Node) GetProxy() *Proxy {
	return n.proxy
}

func (n *Node) Start() {
	go n.startGrpc()
}

func (n *Node) Stop() {
	n.registry.Unregister(context.Background(), &discovery.Node{
		Name: "im",
		Host: "127.0.0.1",
		Port: 9000,
	})
}

func (n *Node) startGrpc() {
	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	defer s.GracefulStop()
	node.RegisterNodeServer(s, &nodeService{node: n})

	r, err := etcd.NewRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	n.registry = r

	err = n.registry.Register(context.Background(), &discovery.Node{
		Name: "im",
		Host: "127.0.0.1",
		Port: 9000,
	}, 10)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
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
