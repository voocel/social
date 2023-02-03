package tcp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"social/pkg/log"
	"social/pkg/network"
	"sync"
	"sync/atomic"
	"time"
)

type Conn struct {
	rw       sync.RWMutex
	cid      int64
	uid      int64
	state    int32
	conn     net.Conn
	sendCh   chan []byte
	msgCh    chan []byte
	done     chan struct{}
	errDone  chan error
	srv      *server
	timer    *time.Timer
	sessId   string
	interval time.Duration
	extraMap map[string]interface{}
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
	_, err := c.conn.Write(msg)
	return err
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
	delete(c.srv.conns, c.conn)
	err := c.conn.Close()
	c.conn = nil
	c.srv.pool.Put(c)

	if c.srv.disconnectHandler != nil {
		c.srv.disconnectHandler(c, err)
	}
	return err
}

func (c *Conn) LocalIP() string {
	addr := c.LocalAddr()

	return ExtractIP(addr)
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
		return "unknown"
	}

	return c.conn.RemoteAddr().String()
}

func (c *Conn) Values() url.Values {
	return nil
}

func (c *Conn) process(ctx context.Context) {
	sess := NewSession(c)
	c.sessId = sess.GetSessionID()
	c.srv.sessions.Store(c.sessId, sess)
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.conn.Close()
		c.srv.sessions.Delete(sess.GetSessionID())
	}()

	go c.readLoop(ctx)
	go c.writeLoop(ctx)

	if c.srv.connectHandler != nil {
		c.srv.connectHandler(c)
	}

	for {
		select {
		case <-c.srv.exitCh:
			return
		case msg := <-c.msgCh:
			c.srv.receiveHandler(c, msg)
		}
	}
}

// readLoop read goroutine
func (c *Conn) readLoop(ctx context.Context) {
	defer c.Close()
	reader := bufio.NewReader(c.conn)
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		default:
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("read err: ", err)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					err = ErrClientClosed
				}
				return
			}
			c.msgCh <- bytes.Trim(buf, "\n")

			v, ok := c.srv.sessions.Load(c.sessId)
			if !ok {
				return
			}
			sess := v.(*Session)
			sess.UpdateTime()
		}
	}
}

// writeLoop write goroutine
func (c *Conn) writeLoop(ctx context.Context) {
	defer c.Close()
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		case msg := <-c.sendCh:
			_, err := c.conn.Write(msg)
			if err != nil {
				log.Errorf("write message err: %v", err)
			}
		case <-c.timer.C:
			//c.SendBytes(Heartbeat, []byte("ping"))
			if c.interval > 0 {
				c.timer.Reset(c.interval)
			}
		}
	}
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

func ExtractIP(addr string) (host string) {
	host, _, _ = net.SplitHostPort(addr)
	return
}
