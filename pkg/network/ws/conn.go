package ws

import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"social/pkg/log"
	"social/pkg/network"
)

type Conn struct {
	rw     sync.RWMutex
	cid    int64
	uid    int64
	state  int32
	conn   *websocket.Conn
	msgCh  chan []byte
	sendCh chan []byte
	exitCh chan struct{}
	server *server
	values url.Values
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

func (c *Conn) Send(msg []byte) error {
	if err := c.checkState(); err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Conn) Push(msg []byte) error {
	if err := c.checkState(); err != nil {
		return err
	}
	c.sendCh <- msg
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
	c.values = nil
	c.server.pool.Put(c)

	if c.server.disconnectHandler != nil {
		fmt.Println("开始断开连接")
		c.server.disconnectHandler(c, err)
	}
	close(c.exitCh)
	return err
}

func (c *Conn) LocalIP() string {
	return ExtractIP(c.LocalAddr())
}

func (c *Conn) LocalAddr() string {
	if err := c.checkState(); err != nil {
		return "unknown"
	}

	return c.conn.LocalAddr().String()
}

func (c *Conn) RemoteIP() string {
	return ExtractIP(c.RemoteAddr())
}

func (c *Conn) RemoteAddr() string {
	if err := c.checkState(); err != nil {
		return ""
	}

	return c.conn.RemoteAddr().String()
}

func (c *Conn) Values() url.Values {
	return c.values
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
	for {
		select {
		case <-c.exitCh:
			return
		default:
			//_, readerMsg, err := c.conn.NextReader()
			_, byteMsg, err := c.conn.ReadMessage()
			if err != nil {
				log.Errorf("read message err: %v", err)
				return
			}
			c.msgCh <- byteMsg
		}
	}
}

func (c *Conn) writeLoop() {
	for {
		select {
		case <-c.exitCh:
			return
		case msg := <-c.sendCh:
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error("client write closed")
				return
			}
		}
	}
}
