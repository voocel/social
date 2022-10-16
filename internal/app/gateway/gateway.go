package gateway

import (
	"context"
	"os"
	"os/signal"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
	"social/pkg/network/ws"
	"syscall"
	"time"
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
	opts     *options
	proxy    *proxy
	protocol message.Protocol
}

func NewGateway(opts ...OptionFunc) *Gateway {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	g := &Gateway{opts: o, protocol: message.NewDefaultProtocol()}
	g.proxy = newProxy(g)

	return g
}

func (g *Gateway) Start() {
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnReceive(g.handleReceive)
	g.opts.server.OnDisconnect(g.handleDisconnect)
}

func (g *Gateway) Stop() {
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("gateway server stop failed: %v", err)
	}
}

func (g *Gateway) handleConnect(conn network.Conn) {

}

func (g *Gateway) handleReceive(conn network.Conn, msg *message.Message, msgType int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	g.proxy.push(ctx, conn.Cid(), conn.Uid(), msg)
}

func (g *Gateway) handleDisconnect(conn network.Conn, err error) {

}
