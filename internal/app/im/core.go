package im

import (
	"fmt"
	"social/internal/node"
)

type core struct {
	proxy *node.Proxy
}

func NewCore(proxy *node.Proxy) *core {
	return &core{proxy: proxy}
}

func (c *core) Init() {
	c.proxy.AddRouteHandler(10, c.connect)
	c.proxy.AddRouteHandler(11, c.chat)
	c.proxy.SetDefaultRouteHandler(c.Default)
}

func (c *core) Default(req node.Request) {
	fmt.Println("im 默认收到: ", req.Buffer)
	data := req.Buffer
	b, ok := data.([]byte)
	if ok {
		fmt.Println(string(b))
	}
}

func (c *core) connect(req node.Request) {
	data := req.Buffer
	b, ok := data.([]byte)
	if ok {
		fmt.Println(string(b))
	}
}

func (c *core) chat(req node.Request) {
	data := req.Buffer
	b, ok := data.([]byte)
	if ok {
		fmt.Println(string(b))
	}
}
