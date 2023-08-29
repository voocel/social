package grpc

import (
	"social/internal/transport"
	"social/internal/transport/grpc/gate"
	"social/internal/transport/grpc/node"
)

type Transporter struct {
	opts *transport.Options
}

func NewTransporter(opts ...transport.Option) *Transporter {
	o := transport.DefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Transporter{opts: o}
}

func (t *Transporter) Options() *transport.Options {
	return t.opts
}

func (t *Transporter) NewGateServer(provider transport.GateProvider) transport.Server {
	return gate.NewServer(provider, t.opts)
}

func (t *Transporter) NewNodeServer(provider transport.NodeProvider) transport.Server {
	return node.NewServer(provider, t.opts)
}

func (t *Transporter) NewGateClient(serviceName string) (transport.GateClient, error) {
	return gate.NewClient(serviceName)
}

func (t *Transporter) NewNodeClient(serviceName string) (transport.NodeClient, error) {
	return node.NewClient(serviceName)
}
