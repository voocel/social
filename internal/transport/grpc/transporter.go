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
	return gate.NewServer(provider, t.opts)
}

func (t *Transporter) NewNodeServer(provider transport.NodeProvider) transport.Server {
	return node.NewServer(provider, t.opts)
}

func (t *Transporter) NewGateClient(addr string) (transport.GateClient, error) {
	return gate.NewClient(addr)
}

func (t *Transporter) NewNodeClient(addr string) (transport.NodeClient, error) {
	//TODO implement me
	panic("implement me")
}
