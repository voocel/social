package gateway

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"social/internal/app/gateway/packet"
	"social/internal/entity"
	"social/internal/route"
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
	mu            sync.Mutex
	proxy         *proxy
	instance      *discovery.Node
	srv           transport.Server
	registry      *etcd.Registry
	protocol      message.Protocol
	sessionGroup  *session.SessionGroup
	done          chan struct{}
	nodeEndpoints map[string]*discovery.Node
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
		mu:            sync.Mutex{},
		sessionGroup:  session.NewSessionGroup(),
		done:          make(chan struct{}),
		protocol:      message.NewDefaultProtocol(),
		nodeEndpoints: make(map[string]*discovery.Node),
	}
	g.proxy = newProxy(g)

	return g
}

func (g *Gateway) Start() {
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
		t := time.NewTimer(time.Second * 5)
		for {
			select {
			case <-g.done:
				return
			case <-t.C:
				g.nodeEndpoints["node"] = g.registry.Query("im")
				t.Reset(time.Second * 5)
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
	uid, err := g.parseUid(conn)
	if err != nil {
		b, _ := packet.Pack(&packet.Message{
			Seq:   0,
			Route: route.Auth,
		})
		conn.Send(b)
		log.Errorf("[Gateway] user connect parse uid err: %v", err)
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()
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
	log.Debugf("[Gateway] receive message route: %v(%v), cid: %v, uid: %v,data: %v",
		route.RouteMap[msg.Route], msg.Route, conn.Cid(), conn.Uid(), string(msg.Buffer))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch msg.Route {
	case route.Heartbeat:
		g.heartbeat(conn)
		return
	}

	err = g.proxy.push(ctx, conn.Cid(), conn.Uid(), msg.Route, msg.Buffer)
	if err != nil {
		log.Errorf("GRPC push to node error: %v", err)
	}
}

func (g *Gateway) handleDisconnect(conn network.Conn, err error) {
	log.Debugf("[Gateway] user connection disconnected: %v, err: %v", conn.RemoteAddr(), err)
	g.sessionGroup.RemoveByUid(conn.Uid())
	g.sessionGroup.RemoveByCid(conn.Cid())
}

// send heartbeat
func (g *Gateway) heartbeat(conn network.Conn) {
	b, _ := packet.Pack(&packet.Message{
		Seq:   0,
		Route: route.Heartbeat,
	})
	conn.Send(b)
}

// 上线后处理离线消息
func (g *Gateway) offlineMessage(uid int64, sess *session.Session) {
	for {
		select {
		case <-g.done:
			return
		case msg := <-entity.MsgCache.Get(uid):
			resp := new(entity.Response)
			if err := sess.Push(resp.Wrap(msg)); err != nil {
				log.Errorf("[Gateway] push offline message to user err: %v", err)
			}
		default:
			return
		}
	}
}

// 从token中解析uid
func (g *Gateway) parseUid(conn network.Conn) (uid int64, err error) {
	token, ok := conn.Values()["token"]
	if !ok {
		err = fmt.Errorf("token non-existent: %v", token)
		return
	}
	var claims *jwt.Claims
	claims, err = jwt.ParseToken(token[0])
	if err != nil {
		err = fmt.Errorf("token parse fail: %v", err)
		return
	}

	return claims.User.ID, nil
}
