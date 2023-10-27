package entity

import (
	"social/protos/pb"
	"sync"
)

var MsgCache = &msgCache{
	receiver: sync.Map{},
}

type msgCache struct {
	rw       sync.RWMutex
	receiver sync.Map
}

func (c *msgCache) Add(uid int64, m *pb.MsgEntity) {
	if v, ok := c.receiver.Load(uid); ok {
		if msgChan := v.(chan *pb.MsgEntity); msgChan != nil {
			msgChan <- m
		}
	} else {
		msgChan := make(chan *pb.MsgEntity, 1024)
		msgChan <- m
		c.receiver.Store(uid, msgChan)
	}
}

func (c *msgCache) Get(uid int64) chan *pb.MsgEntity {
	if v, ok := c.receiver.Load(uid); ok {
		return v.(chan *pb.MsgEntity)
	}
	return make(chan *pb.MsgEntity)
}
