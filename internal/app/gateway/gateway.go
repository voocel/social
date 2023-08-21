package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"

	"social/internal/app/gateway/packet"
	"social/internal/entity"
	"social/internal/session"
	"social/internal/transport"
	"social/pkg/discovery"
	"social/pkg/discovery/etcd"
	"social/pkg/jwt"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
	"social/protos/pb"
)

type Gateway struct {
	opts          *options
	proxy         *proxy
	endpoints     *Endpoint
	srv           transport.Server
	registry      *etcd.Registry
	protocol      message.Protocol
	nodeClient    map[string]pb.NodeClient
	sessionGroup  *session.SessionGroup
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
		sessionGroup:  session.NewSessionGroup(),
		done:          make(chan struct{}),
		nodeClient:    make(map[string]pb.NodeClient),
		protocol:      message.NewDefaultProtocol(),
		nodeEndpoints: make(map[string]*discovery.Node),
	}
	g.proxy = newProxy(g)

	return g
}

func (g *Gateway) Start() {
	// 启动node RPC客户端
	g.startNodeClient("im")
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

	fmt.Println("启动ws")
	fmt.Println(g.opts)
	if err := g.opts.server.Start(); err != nil {
		panic(err)
	}
	fmt.Println("ws ok")

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
	g.srv = g.opts.transporter.NewGateServer(&provider{g})

	r, err := etcd.NewRegistry([]string{viper.GetString("etcd.addr")})
	if err != nil {
		panic(err)
	}

	instance := &discovery.Node{
		Name: "gate-inter-rpc-client",
		Addr: viper.GetString("gaterpc.addr"),
	}
	g.registry = r
	err = r.Register(context.Background(), instance, 60)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := g.srv.Start(); err != nil {
			log.Fatalf("GRPC failed to start: %s", err)
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
	g.nodeClient[serviceName] = pb.NewNodeClient(conn)
}

func (g *Gateway) Stop() {
	g.srv.Stop()
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
	s := session.NewSession(conn)
	g.sessionGroup.UidSession[uid] = s
	g.sessionGroup.CidSession[conn.Cid()] = s
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
