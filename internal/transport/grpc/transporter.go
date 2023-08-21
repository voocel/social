package grpc

import (
	"social/internal/transport"
	"social/internal/transport/grpc/gate"
	"social/internal/transport/grpc/node"
)

type Transporter struct {
	opts *Options
}

func NewTransporter(opts ...Option) *Transporter {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Transporter{opts: o}
}

func (t *Transporter) NewGateServer(provider transport.GateProvider) transport.Server {
	return gate.NewServer(provider, &gate.Options{
		Addr:       t.opts.Server.Addr,
		CertFile:   t.opts.Server.CertFile,
		ServerName: t.opts.Server.KeyFile,
	})
}

func (t *Transporter) NewNodeServer(provider transport.NodeProvider) transport.Server {
	return node.NewServer(provider, &node.Options{
		Addr:       t.opts.Server.Addr,
		CertFile:   t.opts.Server.CertFile,
		ServerName: t.opts.Server.KeyFile,
	})
}

func (t *Transporter) NewGateClient(addr string) (transport.GateClient, error) {
	return gate.NewClient(addr)
}

func (t *Transporter) NewNodeClient(addr string) (transport.NodeClient, error) {
	return node.NewClient(addr)
}
