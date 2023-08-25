package transport

import (
	"context"
	"social/internal/event"
	"social/internal/session"
)

type Server interface {
	// Addr 监听地址
	Addr() string
	// Name 传输类型
	Name() string
	// Start 启动服务器
	Start() error
	// Stop 停止服务器
	Stop() error
}

type GateProvider interface {
	Session(target int64) (*session.Session, error)
	// Bind 绑定用户与网关间的关系
	Bind(ctx context.Context, uid int64) error
	// Unbind 解绑用户与网关间的关系
	Unbind(ctx context.Context, uid int64) error
	// Push 发送消息（异步）
	Push(target int64, msg []byte, msgType ...int) error
	// Broadcast 推送广播消息（异步）
	Broadcast(msg []byte, msgType ...int) (n int, err error)
}

type NodeProvider interface {
	// Trigger 触发事件
	Trigger(event event.Event, gid string, uid int64)
	// Deliver 投递消息
	Deliver(gid, nid string, cid, uid int64, message *Message)
}
