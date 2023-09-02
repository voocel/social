package group

import (
	"encoding/json"
	"social/internal/node"
	"social/internal/router"
	"social/pkg/log"
	"social/protos/pb"
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
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, &msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Default receive data: %v", msg)
	return
}

func (c *core) connect(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, &msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Connect receive data: %v", msg)
	return
}

func (c *core) disconnect(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, &msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Disconnect receive data: %v", msg)
	return
}

func (c *core) message(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, &msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Message receive data: %v", msg)
	return
}
