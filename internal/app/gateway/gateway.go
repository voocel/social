package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"social/config"
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
)

type Gateway struct {
	opts          *options
	proxy         *proxy
	endpoints     *Endpoint
	instance      *discovery.Node
	srv           transport.Server
	registry      *etcd.Registry
	protocol      message.Protocol
	nodeClient    map[string]transport.NodeClient
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
		nodeClient:    make(map[string]transport.NodeClient),
		protocol:      message.NewDefaultProtocol(),
		nodeEndpoints: make(map[string]*discovery.Node),
	}
	g.proxy = newProxy(g)

	return g
}

func (g *Gateway) Start() {
	// 启动node RPC客户端
	for _, v := range config.Conf.Transport.Discovery {
		g.startNodeClient(v.Name)
	}
	// 启动RPC服务端
	g.startRPCServer()
	// 启动网关
	g.startGate()
}

func (g *Gateway) startGate() {
	startupMessage(viper.GetString("gateway.addr"), viper.GetString("gateway.name"))
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
	g.srv = g.opts.transporter.NewGateServer(&provider{g})

	r, err := etcd.NewRegistry([]string{viper.GetString("etcd.addr")})
	if err != nil {
		panic(err)
	}
	g.registry = r

	instance := &discovery.Node{
		Name: g.opts.transporter.Options().Name,
		Addr: g.opts.transporter.Options().Server.Addr,
	}
	g.instance = instance

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

func (g *Gateway) getNodeClient(serviceName string) (transport.NodeClient, error) {
	return g.opts.transporter.NewNodeClient(serviceName)
}

// 启动rpc客户端
func (g *Gateway) startNodeClient(serviceName string) {
	nc, err := g.getNodeClient(serviceName)
	if err != nil {
		log.Fatalf("[Gateway] grpc client connect to node [%s] err: %v", serviceName, err)
	}

	log.Infof("[Gateway] grpc client connect to node [%s] is successful!", serviceName)
	g.nodeClient[serviceName] = nc
}

func (g *Gateway) Stop() {
	if err := g.registry.Unregister(context.Background(), g.instance); err != nil {
		log.Errorf("[%s]gateway registry unregister err: %v", g.instance.Name, err)
	}
	g.srv.Stop()
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("gateway server stop failed: %v", err)
	}
	close(g.done)
	log.Infof("[Gateway] stop and unregister successful: %v", g.instance.Name)
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
	go g.offlineMessage(uid, s)
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

	err = g.proxy.push(ctx, conn.Cid(), conn.Uid(), msg.Buffer, msg.Route)
	if err != nil {
		log.Errorf("GRPC push to node error: %v", err)
	}
}

func (g *Gateway) handleDisconnect(conn network.Conn, err error) {
	log.Debugf("[Gateway] user connection disconnected: %v, err: %v", conn.RemoteAddr(), err)
	g.sessionGroup.RemoveByUid(conn.Uid())
	g.sessionGroup.RemoveByCid(conn.Cid())
}

// 上线后处理离线消息
func (g *Gateway) offlineMessage(uid int64, sess *session.Session) {
	for {
		select {
		case <-g.done:
			return
		case msg := <-entity.MsgCache.Get(uid):
			resp := new(entity.Response)
			if err := sess.Push(resp.Resp(msg)); err != nil {
				log.Errorf("[Gateway] push offline message to user err: %v", err)
			}
		default:
			return
		}
	}
}
