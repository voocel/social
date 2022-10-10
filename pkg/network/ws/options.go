package ws

import (
	"net/http"
)

type OptionFunc func(o *Options)

type CheckOriginFunc func(r *http.Request) bool

type Options struct {
	addr        string          // 监听地址
	maxConnNum  int             // 最大连接数
	certFile    string          // 证书文件
	keyFile     string          // 秘钥文件
	path        string          // 路径，默认为"/"
	checkOrigin CheckOriginFunc // 跨域检测
}

func defaultOptions() *Options {
	return &Options{}
}

func WithTLS(certFile, keyFile string) OptionFunc {
	return func(o *Options) {
		o.certFile = certFile
		o.keyFile = keyFile
	}
}
