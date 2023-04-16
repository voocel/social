package group

import (
	"encoding/json"
	"social/internal/entity"
	"social/internal/node"
	"social/pkg/log"
)

type core struct {
	proxy *node.Proxy
}

func newCore(proxy *node.Proxy) *core {
	return &core{proxy: proxy}
}

func (c *core) Init() {
	c.proxy.AddRouteHandler(10, c.connect)
	c.proxy.AddRouteHandler(11, c.chat)
	c.proxy.SetDefaultRouteHandler(c.Default)
}

func (c *core) Default(req node.Request) error {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return err
	}
	log.Debugf("[group] default receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return nil
}

func (c *core) connect(req node.Request) error {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return err
	}
	log.Debugf("[group] connect receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return nil
}

func (c *core) chat(req node.Request) error {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		return err
	}
	log.Debugf("[group] chat receive cmd: %v, params: %v", arg.Cmd, arg.Params)
	return nil
}
