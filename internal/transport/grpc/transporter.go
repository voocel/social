package grpc

import (
	"social/internal/transport"
	"social/internal/transport/grpc/gate"
)

type Transporter struct {
	opts *options
}

func NewTransporter(opts ...Option) *Transporter {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Transporter{opts: o}
}

func (t *Transporter) NewGateServer() (transport.Server, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Transporter) NewNodeServer() (transport.Server, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Transporter) NewGateClient(addr string) (transport.GateClient, error) {
	return gate.NewClient(addr)
}

func (t *Transporter) NewNodeClient(addr string) (transport.NodeClient, error) {
	//TODO implement me
	panic("implement me")
}
