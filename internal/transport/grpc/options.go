package grpc

import "google.golang.org/grpc"

type Option func(o *options)

type options struct {
	server struct {
		addr       string
		certFile   string
		keyFile    string
		serverOpts []grpc.ServerOption
	}
	client struct {
		certFile   string
		serverName string
	}
}

func defaultOptions() *options {
	opts := &options{}

	return opts
}
