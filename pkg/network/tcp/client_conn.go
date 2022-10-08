package tcp

import (
	"fmt"
	"net"
	"social/pkg/message"
	"social/pkg/network"
	"sync"
	"sync/atomic"
)

type clientConn struct {
	rw     sync.RWMutex
	cid    int64
	uid    int64
	conn   net.Conn
	state  int32
	client *client
	sendCh chan *message.Message
	done   chan struct{}
}

func newClientConn(c *client, conn net.Conn) network.Conn {
	cc := &clientConn{
		cid:    1,
		conn:   conn,
		state:  int32(network.ConnOpened),
		client: c,
		sendCh: make(chan *message.Message, 1024),
		done:   make(chan struct{}),
	}
	if cc.client.connectHandler != nil {
		cc.client.connectHandler(cc)
	}

	return cc
}

func (c *clientConn) readLoop()  {
	for {
		select {
		case <-c.done:
			return
		default:
			msg, err := c.client.protocol.Unpack(c.conn)
			if err != nil {
				fmt.Println("unpack err: ", err)
				c.conn.Close()
				return
			}
			if c.client.receiveHandler != nil {
				c.client.receiveHandler(c, msg, 1)
			}
		}
	}
}

func (c *clientConn) writeLoop()  {

}

func (c *clientConn) Cid() int64 {
	return c.cid
}

func (c *clientConn) Uid() int64 {
	return c.uid
}

func (c *clientConn) Bind(uid int64) {
	c.uid = uid
}

func (c *clientConn) Send(msg []byte, msgType ...int) error {
	if err := c.checkState(); err != nil {
		return err
	}
	_, err := c.conn.Write(msg)
	return err
}

func (c *clientConn) AsyncSend(msg []byte, msgType ...int) error {
	if err := c.checkState(); err != nil {
		return err
	}
	m := message.NewMessage(message.Heartbeat, msg)
	c.sendCh <- m
	return nil
}

func (c *clientConn) State() network.ConnState {
	return network.ConnState(atomic.LoadInt32(&c.state))
}

func (c *clientConn) Close() error {
	if err := c.checkState(); err != nil {
		return err
	}

	close(c.sendCh)
	return c.conn.Close()
}

func (c *clientConn) LocalIP() (string, error) {
	addr, err := c.LocalAddr()
	if err != nil {
		return "", err
	}

	return ExtractIP(addr)
}

func (c *clientConn) LocalAddr() (net.Addr, error) {
	if err := c.checkState(); err != nil {
		return nil, err
	}

	return c.conn.LocalAddr(), nil
}

func (c *clientConn) RemoteIP() (string, error) {
	addr, err := c.RemoteAddr()
	if err != nil {
		return "", err
	}
	return ExtractIP(addr)
}

func (c *clientConn) RemoteAddr() (net.Addr, error) {
	if err := c.checkState(); err != nil {
		return nil, err
	}

	return c.conn.RemoteAddr(), nil
}

func (c *clientConn) checkState() error {
	switch network.ConnState(atomic.LoadInt32(&c.state)) {
	case network.ConnHanged:
		return network.ErrConnectionHanged
	case network.ConnClosed:
		return network.ErrConnectionClosed
	}

	return nil
}