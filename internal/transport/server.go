package transport

import (
	"social/internal/event"
	"social/internal/session"
	"social/protos/pb"
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
	// Push 发送消息（异步）
	Push(req *pb.PushReq) error
	// Multicast 组播消息（异步）
	Multicast(targets []int64, req *pb.Message) (n int64)
	// Broadcast 广播消息（异步）
	Broadcast(req *pb.Message) (n int64)
}

type NodeProvider interface {
	// Trigger 触发事件
	Trigger(event event.Event, gid string, uid int64)
	// Deliver 投递消息
	Deliver(cid, uid int64, message *Message)
}
