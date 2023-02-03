package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"social/pkg/network"
	"sync"
	"sync/atomic"
	"time"
)

type clientConn struct {
	rw       sync.RWMutex
	cid      int64
	uid      int64
	conn     net.Conn
	state    int32
	client   *client
	sendCh   chan []byte
	done     chan struct{}
	timer    *time.Timer
	interval time.Duration
}

func newClientConn(c *client, conn net.Conn) network.Conn {
	cc := &clientConn{
		cid:    1,
		conn:   conn,
		state:  int32(network.ConnOpened),
		client: c,
		sendCh: make(chan []byte, 1024),
		done:   make(chan struct{}),
		timer:  time.NewTimer(10 * time.Second),
	}

	go cc.readLoop()
	go cc.writeLoop()

	if cc.client.connectHandler != nil {
		cc.client.connectHandler(cc)
	}

	return cc
}

func (c *clientConn) readLoop() {
	defer c.Close()
	reader := bufio.NewReader(c.conn)
	for {
		select {
		case <-c.done:
			return
		default:
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					fmt.Println("read err: ", err)
				}
				return
			}
			if c.client.receiveHandler != nil {
				c.client.receiveHandler(c, bytes.Trim(buf, "\n"))
			}
		}
	}
}

func (c *clientConn) writeLoop() {
	defer c.Close()
	for {
		select {
		case <-c.done:
			return
		case msg := <-c.sendCh:
			_, err := c.conn.Write(msg)
			if err != nil {
				fmt.Println("write message err: ", err)
			}
		case <-c.timer.C:
			//c.SendBytes(Heartbeat, []byte("ping"))
			if c.interval > 0 {
				c.timer.Reset(c.interval)
			}
		}
	}
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

func (c *clientConn) Send(msg []byte) error {
	if err := c.checkState(); err != nil {
		return err
	}
	_, err := c.conn.Write(msg)
	return err
}

func (c *clientConn) Push(msg []byte) error {
	if err := c.checkState(); err != nil {
		return err
	}
	c.sendCh <- msg
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

func (c *clientConn) LocalIP() string {
	return ExtractIP(c.LocalAddr())
}

func (c *clientConn) LocalAddr() string {
	if err := c.checkState(); err != nil {
		return "unknown"
	}

	return c.conn.LocalAddr().String()
}

func (c *clientConn) RemoteIP() string {
	return ExtractIP(c.RemoteAddr())
}

func (c *clientConn) RemoteAddr() string {
	if err := c.checkState(); err != nil {
		return ""
	}

	return c.conn.RemoteAddr().String()
}

func (c *clientConn) Values() url.Values {
	return nil
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
