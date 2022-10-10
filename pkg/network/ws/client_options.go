package ws

import "crypto/tls"

type clientOptions struct {
	tlsConf *tls.Config
}

type ClientOptionFunc func(o *clientOptions)

func defaultClientOptions() *clientOptions {
	return &clientOptions{}
}

func WithClientTLS(tls *tls.Config) ClientOptionFunc {
	return func(o *clientOptions) {
		o.tlsConf = tls
	}
}
