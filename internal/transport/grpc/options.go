package grpc

import "google.golang.org/grpc"

const (
	defaultServerAddr = ":7400" // 默认服务器地址
)

type Option func(o *Options)

type Options struct {
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

func defaultOptions() *Options {
	opts := &Options{}
	opts.Server.Addr = defaultServerAddr
	return opts
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
