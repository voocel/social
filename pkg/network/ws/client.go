package ws

import (
	"github.com/gorilla/websocket"
	"social/pkg/message"
	"social/pkg/network"
	"time"
)

type client struct {
	addr     string
	opts     *clientOptions
	dialer   *websocket.Dialer
	protocol message.Protocol

	connectHandler    network.ConnectHandler
	disconnectHandler network.DisconnectHandler
	receiveHandler    network.ReceiveHandler
}

func NewClient(addr string, opts ...ClientOptionFunc) network.Client {
	o := defaultClientOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &client{
		addr: addr,
		opts: o,
		dialer: &websocket.Dialer{
			HandshakeTimeout: time.Second * 5,
		},
		protocol: message.NewDefaultProtocol(),
	}
}

func (c *client) Dial() (network.Conn, error) {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		panic(err)
	}
	return NewWsConn(c, conn), err
}

func (c *client) OnConnect(handler network.ConnectHandler) {
	c.connectHandler = handler
}

func (c *client) OnReceive(handler network.ReceiveHandler) {
	c.receiveHandler = handler
}

func (c *client) OnDisconnect(handler network.DisconnectHandler) {
	c.disconnectHandler = handler
}

func (c *client) Codec() message.Protocol {
	return c.protocol
}


