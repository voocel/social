package im

import (
	"context"
	"encoding/json"
	"social/internal/cmd"
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
	c.proxy.AddRouteHandler(cmd.Connect, c.Connect)
	c.proxy.AddRouteHandler(cmd.Chat, c.Chat)
	c.proxy.SetDefaultRouteHandler(c.Default)
}

func (c *core) Default(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		log.Error(err)
		return
	}
	log.Debugf("[im] default receive cmd: %v, params: %v", arg.Cmd, arg.Params)
}

func (c *core) Connect(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		log.Error(err)
		return
	}
	log.Debugf("[im] connect receive cmd: %v, params: %v", arg.Cmd, arg.Params)
}

func (c *core) Chat(req node.Request) {
	var arg entity.Request
	data := req.Buffer
	if err := json.Unmarshal(data, &arg); err != nil {
		log.Error(err)
		return
	}
	log.Debugf("[im] chat receive cmd: %v, params: %v", arg.Cmd, arg.Params)

	c.proxy.Respond(context.Background(), int64(arg.Params.Receiver), []byte(arg.Params.Content))
}
