package tcp

import (
	"crypto/tls"
	"time"
)

type Options struct {
	logger    Logger
	tlsConf   *tls.Config
	heartbeat time.Duration
}

type Option func(o *Options)

func defaultOptions() *Options {
	return &Options{
		logger:    newLogger(),
		tlsConf:   nil,
		heartbeat: 0,
	}
}

func WithTLS(tls *tls.Config) Option {
	return func(o *Options) {
		o.tlsConf = tls
	}
}

func WithHeartbeat(t time.Duration) Option {
	return func(o *Options) {
		o.heartbeat = t
	}
}

func WithLogger(log Logger) Option {
	return func(o *Options) {
		o.logger = log
		SetLogger(log)
	}
}
