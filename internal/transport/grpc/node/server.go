package node

import (
	"google.golang.org/grpc"
	"net"
	"social/internal/transport"
	tgrpc "social/internal/transport/grpc"
	"social/protos/pb"
)

const name = "grpc"

type server struct {
	addr string
	name string
	srv  *grpc.Server
	lis  net.Listener
}

func NewServer(provider transport.NodeProvider, opts *tgrpc.Options) *server {
	s := grpc.NewServer()
	pb.RegisterNodeServer(s, &nodeService{
		provider: provider,
	})
	return &server{
		addr: opts.Server.Addr,
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

func (s *server) RegisterService(reg func(*grpc.Server)) {
	reg(s.srv)
}

func (s *server) Stop() error {
	s.srv.GracefulStop()
	return s.lis.Close()
}
