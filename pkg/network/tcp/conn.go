package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"social/pkg/log"
	"social/pkg/message"
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
	sendCh   chan *message.Message
	msgCh    chan *message.Message
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

func (c *Conn) Send(msg []byte, msgType ...int) error {
	if err := c.checkState(); err != nil {
		return err
	}
	_, err := c.conn.Write(msg)
	return err
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

	err := c.conn.Close()
	c.conn = nil
	delete(c.srv.conns, c.conn)
	c.srv.pool.Put(c.conn)
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
			c.srv.receiveHandler(c, msg, 1)
		}
	}
}

// readLoop read goroutine
func (c *Conn) readLoop(ctx context.Context) {
	reader := bufio.NewReader(c.conn)
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		default:
			msg, err := c.srv.protocol.Unpack(reader)
			if err != nil {
				fmt.Println("unpack错误: ", err)
				if err == io.EOF {
					err = ErrClientClosed
				} else {
					netOpError, ok := err.(*net.OpError)
					if ok && netOpError.Err.Error() == "use of closed network connection" {
						err = ErrServerClosed
					}
				}
				c.srv.disconnectHandler(c, err)
				c.conn.Close()
				return
			}
			c.msgCh <- msg

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
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		case msg := <-c.sendCh:
			b, err := c.srv.protocol.Pack(msg)
			if err != nil {
				log.Errorf("send message err: %v", err)
			}
			_, err = c.conn.Write(b)
			if err != nil {
				c.srv.disconnectHandler(c, err)
				c.conn.Close()
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

func ExtractIP(addr net.Addr) (host string, err error) {
	host, _, err = net.SplitHostPort(addr.String())
	return
}
