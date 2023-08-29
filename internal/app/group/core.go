package group

import (
	"encoding/json"
	"social/internal/entity"
	"social/internal/node"
	"social/internal/router"
	"social/pkg/log"
)

type core struct {
	proxy *node.Proxy
}

func newCore(proxy *node.Proxy) *core {
	return &core{proxy: proxy}
}

func (c *core) Init() {
	c.proxy.AddRouteHandler(router.Connect, c.connect)
	c.proxy.AddRouteHandler(router.Disconnect, c.disconnect)
	c.proxy.AddRouteHandler(router.GroupMessage, c.message)
	c.proxy.SetDefaultRouteHandler(c.Default)
}

func (c *core) Default(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return
	}
	log.Debugf("[group] default receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return
}

func (c *core) connect(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return
	}
	log.Debugf("[group] connect receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return
}

func (c *core) disconnect(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return
	}
	log.Debugf("[group] disconnect receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return
}

func (c *core) message(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return
	}
	log.Debugf("[group] chat receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return
}
