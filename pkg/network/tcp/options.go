package tcp

import (
	"crypto/tls"
	"time"
)

type Options struct {
	addr      string
	tlsConf   *tls.Config
	heartbeat time.Duration
}

type OptionFunc func(o *Options)

func defaultOptions() *Options {
	return &Options{
		tlsConf:   nil,
		heartbeat: 0,
	}
}

func WithTLS(tls *tls.Config) OptionFunc {
	return func(o *Options) {
		o.tlsConf = tls
	}
}

func WithHeartbeat(t time.Duration) OptionFunc {
	return func(o *Options) {
		o.heartbeat = t
	}
}
