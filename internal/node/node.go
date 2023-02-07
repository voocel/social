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

type RouteHandler func(req Request)

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
	proxy               *Proxy
	opts                *options
	registry            *etcd.Registry
	routes              map[int32]routeEntity
	instance            *discovery.Node
	srv                 *grpc.Server
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
	n.routes = make(map[int32]routeEntity)
	n.instance = instance
	return n
}

func (n *Node) GetProxy() *Proxy {
	return n.proxy
}

func (n *Node) Start() {
	go n.startRPCServer()
}

func (n *Node) Stop() {
	if err := n.registry.Unregister(context.Background(), n.instance); err != nil {
		log.Errorf("[%s]registry unregister err: %v", n.instance.Name, err)
	}
	n.srv.GracefulStop()
}

func (n *Node) startRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.instance.Host, n.instance.Port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	n.srv = s
	node.RegisterNodeServer(s, &nodeService{node: n})

	r, err := etcd.NewRegistry([]string{viper.GetString("etcd.addr")})
	if err != nil {
		panic(err)
	}
	n.registry = r

	err = n.registry.Register(context.Background(), n.instance, 60)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	log.Infof("[%s]node stop success", n.instance.Name)
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
