package gateway

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"social/internal/app/gateway/packet"
	"social/internal/entity"
	"social/pkg/discovery"
	"social/pkg/discovery/etcd"
	"social/pkg/jwt"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
	"social/pkg/network/ws"
	"social/protos/gate"
)

func Run() *Gateway {
	srv := ws.NewServer(":8800")
	gateway := NewGateway(WithServer(srv))
	gateway.Start()
	return gateway
}

type Gateway struct {
	opts          *options
	proxy         *proxy
	endpoints     *Endpoint
	srv           *grpc.Server
	registry      *etcd.Registry
	protocol      message.Protocol
	nodeConns     map[string]*grpc.ClientConn
	sessionGroup  *SessionGroup
	instance      *discovery.Node
	done          chan struct{}
	nodeEndpoints map[string]*discovery.Node
}

type Endpoint struct {
	gateEndpoints map[string]string
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		gateEndpoints: make(map[string]string),
	}
}

func (e *Endpoint) add(id, addr string) {
	e.gateEndpoints[id] = addr
}

func NewGateway(opts ...OptionFunc) *Gateway {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	if o.id == "" {
		o.id = uuid.New().String()
	}
	g := &Gateway{
		opts:          o,
		endpoints:     NewEndpoint(),
		sessionGroup:  NewSessionGroup(),
		done:          make(chan struct{}),
		nodeConns:     make(map[string]*grpc.ClientConn),
		protocol:      message.NewDefaultProtocol(),
		nodeEndpoints: make(map[string]*discovery.Node),
		instance: &discovery.Node{
			Name: "gate-inter-rpc-client",
			Host: viper.GetString("gaterpc.host"),
			Port: viper.GetInt("gaterpc.port"),
		},
	}
	g.proxy = newProxy(g)

	return g
}

func (g *Gateway) Start() {
	// 启动RPC客户端
	g.startNodeClient("im")
	g.proxy.newNodeClient("im")
	// 启动RPC服务端
	g.startRPCServer()
	// 启动网关
	g.startGate()
}

func (g *Gateway) startGate() {
	startupMessage(":8800", "Gateway")
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnReceive(g.handleReceive)
	g.opts.server.OnDisconnect(g.handleDisconnect)

	if err := g.opts.server.Start(); err != nil {
		panic(err)
	}

	// get node service instance
	go func() {
		t := time.NewTimer(time.Second * 10)
		for {
			select {
			case <-g.done:
				return
			case <-t.C:
				g.nodeEndpoints["node"] = g.registry.Query("node")
			}
		}
	}()
}

// 启动rpc服务端
func (g *Gateway) startRPCServer() {
	lis, err := net.Listen("tcp", ":7400")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s := grpc.NewServer()
	g.srv = s
	//s.RegisterService(&gate.Gate_ServiceDesc, &endpoint{})
	gate.RegisterGateServer(s, &endpoint{sessionGroup: g.sessionGroup})

	r, err := etcd.NewRegistry([]string{viper.GetString("etcd.addr")})
	if err != nil {
		panic(err)
	}
	g.registry = r
	err = g.registry.Register(context.Background(), g.instance, 60)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer lis.Close()
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
		log.Infof("gateway GRPC server stop success")
	}()
}

// 启动rpc客户端
func (g *Gateway) startNodeClient(serviceName string) {
	reg, err := etcd.NewResolver([]string{viper.GetString("etcd.addr")}, serviceName)
	if err != nil {
		panic(err)
	}
	resolver.Register(reg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infof("[Gateway] grpc client trying to connect to node [%s]...", serviceName)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", reg.Scheme(), serviceName), grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithBlock())
	if err != nil {
		log.Warnf("[Gateway] the node[%s] grpc server not ready yet: %v", serviceName, err)
		return
	}

	log.Infof("[Gateway] grpc client connect to [%s] is successful!", serviceName)
	g.nodeConns[serviceName] = conn
}

func (g *Gateway) Stop() {
	g.srv.GracefulStop()
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("gateway server stop failed: %v", err)
	}
	close(g.done)
}

func (g *Gateway) handleConnect(conn network.Conn) {
	log.Debugf("[Gateway] user connect successful: %v", conn.RemoteAddr())
	resp := new(entity.Response)
	token, ok := conn.Values()["token"]
	if !ok {
		conn.Send(resp.ErrResp("token non-existent"))
		conn.Close()
		log.Errorf("token non-existent: %v", token)
		return
	}
	claims, err := jwt.ParseToken(token[0])
	if err != nil {
		conn.Send(resp.ErrResp("token parse fail: " + err.Error()))
		conn.Close()
		log.Errorf("token parse fail: %v", err)
		return
	}
	uid := claims.User.ID
	s := newSession(conn)
	g.sessionGroup.uidSession[uid] = s
	g.sessionGroup.cidSession[conn.Cid()] = s
	conn.Bind(uid)
}

func (g *Gateway) handleReceive(conn network.Conn, data []byte) {
	msg, err := packet.Unpack(data)
	if err != nil {
		log.Errorf("unpack data to struct failed: %v", err)
		return
	}
	log.Debugf("[Gateway] receive message route: %v, cid: %v, uid: %v,data: %v", msg.Route, conn.Cid(), conn.Uid(), string(msg.Buffer))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = g.proxy.push(ctx, conn.Cid(), conn.Uid(), msg.Buffer, msg.Route)
	if err != nil {
		log.Errorf("GRPC push to node error: %v", err)
	}
}

func (g *Gateway) handleDisconnect(conn network.Conn, err error) {
	log.Debugf("[Gateway] user connection disconnected: %v, err: %v", conn.RemoteAddr(), err)
	g.sessionGroup.RemoveByUid(conn.Uid())
	g.sessionGroup.RemoveByCid(conn.Cid())
}
