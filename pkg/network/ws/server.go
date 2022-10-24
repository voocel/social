package ws

import (
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
)

type server struct {
	opts     *Options
	listener net.Listener
	pool     sync.Pool
	upgrade  websocket.Upgrader
	protocol message.Protocol
	conns    map[*websocket.Conn]*Conn

	startHandler      network.StartHandler
	stopHandler       network.CloseHandler
	connectHandler    network.ConnectHandler
	disconnectHandler network.DisconnectHandler
	receiveHandler    network.ReceiveHandler
}

func NewServer(addr string, opts ...OptionFunc) network.Server {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	o.addr = addr
	return &server{
		opts:     o,
		listener: nil,
		upgrade: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		pool: sync.Pool{
			New: func() interface{} { return &Conn{} },
		},
		conns:    make(map[*websocket.Conn]*Conn),
		protocol: message.NewDefaultProtocol(),
	}
}

func (s *server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", s.opts.addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return err
	}
	s.listener = l

	if s.startHandler != nil {
		s.startHandler()
	}

	go s.serve()

	return err
}

func (s *server) serve() {
	var err error
	http.HandleFunc("/ws", s.wsHandle)
	if s.opts.certFile != "" && s.opts.keyFile != "" {
		err = http.ServeTLS(s.listener, nil, s.opts.certFile, s.opts.keyFile)
	} else {
		err = http.Serve(s.listener, nil)
	}

	if err != nil {
		log.Fatalf("websocket server start error: %v", err)
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
	return "websocket"
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

func (s *server) wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	log.Infof("连接成功：%v\n", conn.RemoteAddr())

	c := s.pool.Get().(*Conn)
	c.cid++
	c.conn = conn
	c.msgCh = make(chan *message.Message, 1024)
	c.sendCh = make(chan *message.Message, 1024)
	c.exitCh = make(chan struct{})
	c.server = s

	if s.connectHandler != nil {
		s.connectHandler(c)
	}
	s.conns[conn] = c

	go c.readLoop()
	go c.writeLoop()

	for {
		select {
		case <-c.exitCh:
			return
		case msg := <-c.msgCh:
			if s.receiveHandler != nil {
				s.receiveHandler(c, msg, 1)
			}
		}
	}
}