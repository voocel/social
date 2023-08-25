package entity

import (
	"fmt"
	"sync"
)

var MsgCache = &msgCache{
	receiver: sync.Map{},
}

type msgCache struct {
	rw       sync.RWMutex
	receiver sync.Map
}

func (c *msgCache) Add(uid int64, m *Message) {
	if v, ok := c.receiver.Load(uid); ok {
		if msgChan := v.(chan *Message); msgChan != nil {
			fmt.Println("记录离线消息")
			msgChan <- m
		}
	} else {
		msgChan := make(chan *Message, 1024)
		msgChan <- m
		c.receiver.Store(uid, msgChan)
	}
}

func (c *msgCache) Get(uid int64) chan *Message {
	if v, ok := c.receiver.Load(uid); ok {
		return v.(chan *Message)
	}
	return make(chan *Message)
}
