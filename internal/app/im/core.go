package im

import (
	"context"
	"encoding/json"
	"social/ent"
	"social/internal/node"
	"social/internal/route"
	"social/internal/usecase"
	"social/internal/usecase/repo"
	"social/pkg/log"
	"social/protos/pb"
	"time"
)

type core struct {
	proxy       *node.Proxy
	userUseCase *usecase.UserUseCase
}

func newCore(proxy *node.Proxy, entClient *ent.Client) *core {
	return &core{proxy: proxy, userUseCase: usecase.NewUserUseCase(repo.NewUserRepo(entClient))}
}

func (c *core) Init() {
	c.proxy.AddRouteHandler(route.Connect, c.connect)
	c.proxy.AddRouteHandler(route.Disconnect, c.disconnect)
	c.proxy.AddRouteHandler(route.Message, c.message)
	c.proxy.SetDefaultRouteHandler(c.Default)
}

func (c *core) Default(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[IM]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[IM]Default receive data: %v", msg)
	return
}

func (c *core) connect(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[IM]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[IM]Connect receive data: %v", msg)
	return
}

func (c *core) disconnect(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[IM]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[IM]Disconnect receive data: %v", msg)
	return
}

func (c *core) message(req node.Request) {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[IM]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[IM]Message receive data: %v", msg)

	user, err := c.userUseCase.GetUserById(context.Background(), msg.Sender.Id)
	if err != nil {
		log.Errorf("[IM]GetUserById err: %v", err)
		return
	}
	msg.Sender.Nickname = user.Nickname
	msg.Sender.Avatar = user.Avatar
	msg.Timestamp = time.Now().Unix()

	err = req.Respond(context.Background(), msg.Receiver.Id, msg)
	if err != nil {
		log.Errorf("[IM]Respond message err: %v", err)
	}
}
