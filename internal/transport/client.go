package transport

import (
	"context"
	"social/pkg/event"
)

type NodeClient interface {
	// Trigger 触发事件
	Trigger(ctx context.Context, event event.Event, gid string, uid int64) (miss bool, err error)
	// Deliver 投递消息
	Deliver(ctx context.Context, gid, nid string, cid, uid int64, message *Message) (miss bool, err error)
}

type GateClient interface {
	// Bind 绑定用户与连接
	Bind(ctx context.Context, cid, uid int64) (err error)
	// Unbind 解绑用户与连接
	Unbind(ctx context.Context, uid int64) (err error)
	// GetIP 获取客户端IP
	GetIP(ctx context.Context, target int64) (ip string, err error)
	// Disconnect 断开连接
	Disconnect(ctx context.Context, target int64, isForce bool) (err error)
	// Push 推送消息
	Push(ctx context.Context, target int64, message *Message) (err error)
	// Multicast 推送组播消息
	Multicast(ctx context.Context, targets []int64, message *Message) (total int64, err error)
	// Broadcast 推送广播消息
	Broadcast(ctx context.Context, message *Message) (total int64, err error)
}

type Message struct {
	Seq    int32  // 序列号
	Route  int32  // 路由
	Buffer []byte // 消息内容
}
