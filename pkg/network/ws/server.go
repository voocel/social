package ws

import (
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"social/pkg/log"
	"social/pkg/message"
	"social/pkg/network"
)

func DefaultErrorWriter(rw http.ResponseWriter, req *http.Request, code int, err error) {
	rw.WriteHeader(code)
	rw.Write([]byte(err.Error()))
}

type server struct {
	cid      int64
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
			ReadBufferSize:   4096,
			WriteBufferSize:  4096,
			HandshakeTimeout: time.Second * 5,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Error: DefaultErrorWriter,
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
		err := conn.Close()
		if err != nil {
			log.Errorf("closed connection err: %v", err)
		}
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

	s.cid++
	c := s.pool.Get().(*Conn)
	c.cid = s.cid
	c.conn = conn
	c.values = r.URL.Query()
	c.msgCh = make(chan []byte, 1024)
	c.sendCh = make(chan []byte, 1024)
	c.exitCh = make(chan struct{})
	c.server = s
	c.rateLimiter = rate.NewLimiter(rate.Limit(10), 10)

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
				s.receiveHandler(c, msg)
			}
		}
	}
}
