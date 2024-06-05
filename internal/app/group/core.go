package group

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
	proxy     *node.Proxy
	gUseCase  *usecase.GroupUseCase
	gmUseCase *usecase.GroupMemberUseCase
}

func newCore(proxy *node.Proxy, entClient *ent.Client) *core {
	return &core{
		proxy:     proxy,
		gmUseCase: usecase.NewGroupMemberUseCase(repo.NewGroupMemberRepo(entClient)),
		gUseCase:  usecase.NewGroupUseCase(repo.NewGroupRepo(entClient)),
	}
}

func (c *core) Init() {
	c.proxy.AddRouteHandler(route.Connect, c.connect)
	c.proxy.AddRouteHandler(route.Disconnect, c.disconnect)
	c.proxy.AddRouteHandler(route.GroupMessage, c.message)
	c.proxy.SetDefaultRouteHandler(c.Default)
}

func (c *core) Default(req node.Request) {
	var msg = new(pb.MsgEntity)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[Group]Default unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Default receive data: %v", msg)
	return
}

func (c *core) connect(req node.Request) {
	var msg = new(pb.MsgEntity)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Connect receive data: %v", msg)
	return
}

func (c *core) disconnect(req node.Request) {
	var msg = new(pb.MsgEntity)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Disconnect receive data: %v", msg)
	return
}

func (c *core) message(req node.Request) {
	var msg = new(pb.MsgEntity)
	if err := json.Unmarshal(req.Buffer, msg); err != nil {
		log.Errorf("[Group]Unmarshal message err: %v", err)
		return
	}
	log.Debugf("[Group]Message receive data: %v", msg)

	uids, err := c.gmUseCase.GetGroupMemberUser(context.Background(), msg.Receiver.Id)
	if err != nil {
		log.Errorf("[Group]GetGroupMemberUser err: %v", err)
		return
	}

	group, err := c.gUseCase.GetGroupById(context.Background(), msg.Receiver.Id)
	if err != nil {
		log.Errorf("[Group]GetGroupById err: %v", err)
		return
	}

	msg.Receiver.Avatar = group.Avatar
	msg.Receiver.Name = group.Name
	msg.Timestamp = time.Now().Unix()

	err = req.Multicast(context.Background(), uids, msg)
	if err != nil {
		log.Errorf("[Group]Respond Multicast message err: %v", err)
	}
	return
}
