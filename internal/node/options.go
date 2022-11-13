package node

import (
	"social/pkg/discovery"
	"social/pkg/network"
)

type options struct {
	id        string
	name      string
	server    network.Server
	discovery discovery.Discovery
}

type OptionFunc func(o *options)

func defaultOptions() *options {
	return &options{}
}

func WithID(id string) OptionFunc { return func(o *options) { o.id = id } }

func WithName(name string) OptionFunc { return func(o *options) { o.name = name } }

func WithServer(s network.Server) OptionFunc { return func(o *options) { o.server = s } }
