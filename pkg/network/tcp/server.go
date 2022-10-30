package tcp

import (
	"context"
	"errors"
	"net"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
	"sync"
	"time"
)

var (
	ErrServerClosed = errors.New("server closed idle connection")
	ErrClientClosed = errors.New("client closed")
)

const (
	IdleTime = 60
)

type server struct {
	cid      int64
	opts     *Options
	listener net.Listener
	sessions *sync.Map
	pool     sync.Pool
	conns    map[net.Conn]*Conn
	protocol message.Protocol
	exitCh   chan struct{}

	startHandler      network.StartHandler
	stopHandler       network.CloseHandler
	connectHandler    network.ConnectHandler
	disconnectHandler network.DisconnectHandler
	receiveHandler    network.ReceiveHandler
}

func NewServer(addr string, opts ...OptionFunc) *server {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	o.addr = addr
	return &server{
		opts:     o,
		sessions: &sync.Map{},
		exitCh:   make(chan struct{}),
		conns:    make(map[net.Conn]*Conn),
		protocol: message.NewDefaultProtocol(),
		pool: sync.Pool{
			New: func() interface{} { return &Conn{} },
		},
	}
}

func (s *server) Start() error {
	listener, err := net.Listen("tcp", s.opts.addr)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	if s.startHandler != nil {
		s.startHandler()
	}

	go s.serve()

	return nil
}

func (s *server) serve() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.heartbeat()
	s.accept(ctx)
}

func (s *server) accept(ctx context.Context) {
	for {
		select {
		case <-s.exitCh:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("accept connection err: %v", err)
			continue
		}

		s.cid++
		cc := s.pool.Get().(*Conn)
		cc.cid = s.cid
		cc.conn = conn
		cc.timer = time.NewTimer(2 * time.Second)
		cc.msgCh = make(chan *message.Message, 1024)
		cc.sendCh = make(chan *message.Message, 1024)
		cc.extraMap = make(map[string]interface{})
		cc.srv = s

		if s.connectHandler != nil {
			s.connectHandler(cc)
		}
		s.conns[conn] = cc
		go cc.process(ctx)
	}
}

func (s *server) heartbeat() {
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

func (s *server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	for _, conn := range s.conns {
		conn.Close()
	}
	if s.stopHandler != nil {
		s.stopHandler()
	}

	return nil
}

func (s *server) Protocol() string {
	return "tcp"
}

func (s *server) OnStart(handler network.StartHandler) {
	s.startHandler = handler
}

func (s *server) OnStop(handler network.CloseHandler) {
	s.stopHandler = handler
}

func (s *server) OnConnect(handler network.ConnectHandler) {
	s.connectHandler = handler
}

func (s *server) OnReceive(handler network.ReceiveHandler) {
	s.receiveHandler = handler
}

func (s *server) OnDisconnect(handler network.DisconnectHandler) {
	s.disconnectHandler = handler
}
