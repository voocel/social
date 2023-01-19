package gateway

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"social/internal/app/gateway/packet"
	"social/pkg/discovery/etcd"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
	"social/pkg/network/ws"
)

func Run() {
	srv := ws.NewServer(":8800")
	gate := NewGateway(WithServer(srv))
	gate.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Sync()
			gate.Stop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

type Gateway struct {
	opts      *options
	proxy     *proxy
	sessions  *sync.Map
	protocol  message.Protocol
	nodeConns map[string]*grpc.ClientConn
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
		opts:      o,
		sessions:  &sync.Map{},
		nodeConns: make(map[string]*grpc.ClientConn),
		protocol:  message.NewDefaultProtocol(),
	}
	g.proxy = newProxy(g)

	return g
}

func (g *Gateway) Start() {
	g.newNodeClient("im")
	g.proxy.newNodeClient("im")
	g.startGate()
}

func (g *Gateway) startGate() {
	startupMessage(":9000", false, "gateway")
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnReceive(g.handleReceive)
	g.opts.server.OnDisconnect(g.handleDisconnect)

	if err := g.opts.server.Start(); err != nil {
		panic(err)
	}
}

func (g *Gateway) newNodeClient(serviceName string) {
	reg, err := etcd.NewResolver([]string{viper.GetString("etcd.addr")}, serviceName)
	if err != nil {
		panic(err)
	}
	resolver.Register(reg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", reg.Scheme(), serviceName), grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	g.nodeConns[serviceName] = conn
}

func (g *Gateway) Stop() {
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("gateway server stop failed: %v", err)
	}
}

func (g *Gateway) handleConnect(conn network.Conn) {
	fmt.Println("[gateway]连接成功: ", conn.RemoteAddr().String())
	s := newSession(conn)
	g.sessions.Store(conn.Uid(), s)
}

func (g *Gateway) handleReceive(conn network.Conn, data []byte, msgType int) {
	msg, err := packet.Unpack(data)
	if err != nil {
		log.Errorf("unpack data to struct failed: %v", err)
		return
	}
	fmt.Printf("Gateway 收到消息: data: %v, route: %v\n", string(msg.Buffer), msg.Route)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Printf("路由: %v, id: %v, 用户id: %v\n", msg.Route, conn.Cid(), conn.Uid())

	payload, err := g.proxy.push(ctx, conn.Cid(), conn.Uid(), msg.Buffer, msg.Route)
	if err != nil {
		conn.Send([]byte("gateway send error"))
	} else {
		conn.Send(payload)
	}
}

func (g *Gateway) handleDisconnect(conn network.Conn, err error) {
	fmt.Println(conn.RemoteAddr())
	//log.Infof("[Gateway] connection closed: %v, err: %v", conn.RemoteAddr().String(), err.Error())
	g.sessions.Delete(conn.Uid())
}
