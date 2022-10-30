package ws

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
)

type Conn struct {
	rw     sync.RWMutex
	cid    int64
	uid    int64
	state  int32
	conn   *websocket.Conn
	msgCh  chan *message.Message
	sendCh chan *message.Message
	exitCh chan struct{}
	server *server
}

func (c *Conn) Cid() int64 {
	return c.cid
}

func (c *Conn) Uid() int64 {
	return c.uid
}

func (c *Conn) Bind(uid int64) {
	c.uid = uid
}

func (c *Conn) Send(msg []byte, msgType ...int) error {
	if err := c.checkState(); err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Conn) AsyncSend(msg []byte, msgType ...int) error {
	if err := c.checkState(); err != nil {
		return err
	}
	m := message.NewMessage(message.Heartbeat, msg)
	c.sendCh <- m
	return nil
}

func (c *Conn) State() network.ConnState {
	return network.ConnState(atomic.LoadInt32(&c.state))
}

func (c *Conn) Close() error {
	if err := c.checkState(); err != nil {
		return err
	}

	atomic.StoreInt32(&c.state, int32(network.ConnClosed))
	delete(c.server.conns, c.conn)
	err := c.conn.Close()
	c.conn = nil
	c.server.pool.Put(c)

	if c.server.disconnectHandler != nil {
		c.server.disconnectHandler(c, err)
	}
	return err
}

func (c *Conn) LocalIP() (string, error) {
	addr, err := c.LocalAddr()
	if err != nil {
		return "", err
	}

	return ExtractIP(addr)
}

func (c *Conn) LocalAddr() (net.Addr, error) {
	if err := c.checkState(); err != nil {
		return nil, err
	}

	return c.conn.LocalAddr(), nil
}

func (c *Conn) RemoteIP() (string, error) {
	addr, err := c.RemoteAddr()
	if err != nil {
		return "", err
	}
	return ExtractIP(addr)
}

func (c *Conn) RemoteAddr() (net.Addr, error) {
	if err := c.checkState(); err != nil {
		return nil, err
	}

	return c.conn.RemoteAddr(), nil
}

func (c *Conn) checkState() error {
	switch network.ConnState(atomic.LoadInt32(&c.state)) {
	case network.ConnHanged:
		return network.ErrConnectionHanged
	case network.ConnClosed:
		return network.ErrConnectionClosed
	}
	return nil
}

func (c *Conn) readLoop() {
	defer c.Close()
	for {
		select {
		case <-c.exitCh:
			return
		default:
			_, readerMsg, err := c.conn.NextReader()
			if err != nil {
				log.Errorf("read message err: %v", err)
				return
			}
			msg, err := c.server.protocol.Unpack(readerMsg)
			if err != nil {
				log.Errorf("unpack message err: ", err)
				return
			}
			c.msgCh <- msg
		}
	}
}

func (c *Conn) writeLoop() {
	for {
		select {
		case <-c.exitCh:
			return
		case msg := <-c.sendCh:
			byteMsg, err := c.server.protocol.Pack(msg)
			if err != nil {
				log.Errorf("pack message err: ", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, byteMsg); err != nil {
				log.Error("client write closed")
				return
			}
		}
	}
}
