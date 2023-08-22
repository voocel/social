package transport

import (
	"google.golang.org/grpc"
)

const (
	defaultServerAddr = ":7000"
)

type Option func(o *Options)

type Options struct {
	Name   string
	Server struct {
		Addr       string
		CertFile   string
		KeyFile    string
		ServerOpts []grpc.ServerOption
	}
	Client struct {
		certFile   string
		serverName string
	}
}

func DefaultOptions() *Options {
	opts := &Options{}
	opts.Server.Addr = defaultServerAddr
	return opts
}

func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Server.Addr = addr
	}
}

func WithCert(key, cert string) Option {
	return func(o *Options) {
		o.Server.KeyFile = key
		o.Server.CertFile = cert
	}
}
