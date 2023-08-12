package gate

import (
	"google.golang.org/grpc"
	"net"
)

const name = "grpc"

type server struct {
	addr string
	name string
	srv  *grpc.Server
	lis  net.Listener
}

func NewServer(addr string) *server {
	s := grpc.NewServer()
	return &server{
		addr: addr,
		srv:  s,
		name: name,
	}
}

func (s *server) Addr() string {
	return s.addr
}

func (s *server) Name() string {
	return s.name
}

func (s *server) Start() error {
	var err error
	s.lis, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.srv.Serve(s.lis)
}

func (s *server) Stop() error {
	s.srv.Stop()
	return s.lis.Close()
}
