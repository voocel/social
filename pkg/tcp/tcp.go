package tcp

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	// ErrServerClosed occurs when heartbeat detection the connection is not active
	ErrServerClosed = errors.New("server closed idle connection")
	// ErrClientClosed occurs when client close the connection
	ErrClientClosed = errors.New("client closed")
)

const (
	// IdleTime If no data is sent to the server within 60 seconds, the connection will be forcibly closed
	IdleTime = 60
)

// Server defines parameters for running an TCP network
type Server struct {
	addr      string
	opt       *Options
	exitCh    chan struct{}
	sessions  *sync.Map
	onConnect func(c *Conn)
	onMessage func(c *Conn, msg *Message)
	onClose   func(c *Conn, err error)
}

// NewServer creates a new tcp network connection using the given net connection.
func NewServer(addr string, opts ...Option) *Server {
	serv := &Server{
		addr:      addr,
		exitCh:    make(chan struct{}),
		sessions:  &sync.Map{},
		onConnect: func(c *Conn) {},
		onMessage: func(c *Conn, msg *Message) {},
		onClose:   func(c *Conn, err error) {},
	}

	d := defaultOptions()
	for _, opt := range opts {
		opt(d)
	}
	serv.opt = d
	return serv
}

// Start create a TCP network listener to accept client connection
func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	Flog.Infof("TCP server start successfully! %v", s.addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Heartbeat()
	go s.Accept(ctx, listener)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-ch:
			s.Stop()
			return
		}
	}
}

func (s *Server) Accept(ctx context.Context, listener net.Listener) {
	for {
		select {
		case <-s.exitCh:
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			Flog.Errorf("accept connection err: %v", err)
			continue
		}
		c := Conn{
			srv:      s,
			conn:     conn,
			timer:    time.NewTimer(2 * time.Second),
			clientIP: conn.RemoteAddr(),
			protocol: NewDefaultProtocol(),
			msgCh:    make(chan *Message, 1024),
			sendCh:   make(chan *Message, 1024),
			errDone:  make(chan error, 1),
			extraMap: map[string]interface{}{},
		}
		go c.process(ctx)
	}
}

// Heartbeat heartbeat detection
func (s *Server) Heartbeat() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			s.sessions.Range(func(key, value interface{}) bool {
				sess, ok := value.(*Session)
				if !ok {
					return true
				}
				if time.Now().Unix()-sess.lastTime > IdleTime {
					sess.GetConn().Close()
					s.sessions.Delete(key)
				}
				return true
			})
		}
	}
}

// Stop Close server
func (s *Server) Stop() {
	close(s.exitCh)
	s.sessions.Range(func(key, value interface{}) bool {
		sess := value.(*Session)
		sess.GetConn().Close()
		return true
	})
}

// OnConnect connect callbacks on a connection
func (s *Server) OnConnect(callback func(c *Conn)) {
	s.onConnect = callback
}

// OnMessage receive callbacks on a connection
func (s *Server) OnMessage(callback func(c *Conn, msg *Message)) {
	s.onMessage = callback
}

// OnClose close callbacks on a connection, it will be called before closing
func (s *Server) OnClose(callback func(c *Conn, err error)) {
	s.onClose = callback
}

// Conn defines parameters for accept an client
type Conn struct {
	srv      *Server
	conn     net.Conn
	clientIP net.Addr
	protocol Protocol
	sessId   string
	timer    *time.Timer
	timeout  time.Duration
	interval time.Duration
	sendCh   chan *Message
	msgCh    chan *Message
	errDone  chan error
	extraMap map[string]interface{}
	sync.RWMutex
}

// process client connection
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

	c.srv.onConnect(c)
	for {
		select {
		case <-c.srv.exitCh:
			return
		case err := <-c.errDone:
			c.srv.onClose(c, err)
			return
		case msg := <-c.msgCh:
			c.srv.onMessage(c, msg)
		}
	}
}

// readLoop read goroutine
func (c *Conn) readLoop(ctx context.Context) {
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		default:
			reader := bufio.NewReader(c.conn)
			msg, err := c.protocol.Unpack(reader)
			if err != nil {
				if err == io.EOF {
					err = ErrClientClosed
				} else {
					netOpError, ok := err.(*net.OpError)
					if ok && netOpError.Err.Error() == "use of closed network connection" {
						err = ErrServerClosed
					}
				}
				c.errDone <- err
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
			if err := c.writeMessage(msg); err != nil {
				Flog.Errorf("send message err: %v", err)
			}
		case <-c.timer.C:
			//c.SendBytes(Heartbeat, []byte("ping"))
			if c.interval > 0 {
				c.timer.Reset(c.interval)
			}
		}
	}
}

// write Message to client connection
func (c *Conn) writeMessage(msg *Message) error {
	m, err := c.protocol.Pack(msg)
	if err != nil {
		return err
	}
	return c.writeBytes(m)
}

// writeBytes send bytes message to client connection
func (c *Conn) writeBytes(b []byte) error {
	_, err := c.conn.Write(b)
	if err != nil {
		c.conn.Close()
		c.srv.onClose(c, err)
	}
	return err
}

// GetRawConn get the raw net.Conn from the client connection
func (c *Conn) GetRawConn() net.Conn {
	return c.conn
}

// SetExtraMap set the extra data
func (c *Conn) SetExtraMap(k string, v interface{}) {
	c.Lock()
	defer c.Unlock()
	c.extraMap[k] = v
}

// GetExtraMap get the extra data
func (c *Conn) GetExtraMap(k string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.extraMap[k]
}

// GetClientIP get client IP
func (c *Conn) GetClientIP() net.Addr {
	return c.clientIP
}

// SendMessage send message into channel
func (c *Conn) SendMessage(msg *Message) {
	c.sendCh <- msg
}

// SendBytes send bytes
func (c *Conn) SendBytes(cmd CMD, b []byte) {
	msg := NewMessage(cmd, b)
	c.SendMessage(msg)
}

// SendSingle send message to single
func (c *Conn) SendSingle(sid string, msg *Message) {
	v, ok := c.srv.sessions.Load(sid)
	if ok {
		sess := v.(*Session)
		sess.GetConn().SendMessage(msg)
	}
}

// SendAll send message to all
func (c *Conn) SendAll(msg *Message) {
	c.srv.sessions.Range(func(key, value interface{}) bool {
		sess := value.(*Session)
		sess.GetConn().SendMessage(msg)
		return true
	})
}

// Close the client connection
func (c *Conn) Close() {
	c.conn.Close()
	c.errDone <- ErrServerClosed
}
