package gateway

import (
	"social/pkg/discovery"
	"social/pkg/network"
)

const defaultName = "social-gateway"

type options struct {
	id        string
	name      string
	server    network.Server
	discovery discovery.Discovery
}

type OptionFunc func(o *options)

func defaultOptions() *options {
	return &options{
		name: defaultName,
	}
}

func WithID(id string) OptionFunc { return func(o *options) { o.id = id } }

func WithName(name string) OptionFunc { return func(o *options) { o.name = name } }

func WithServer(s network.Server) OptionFunc { return func(o *options) { o.server = s } }
