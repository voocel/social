package network

import (
	"errors"
	"net"
)

const (
	ConnOpened ConnState = iota + 1 // 连接打开
	ConnHanged                      // 连接挂起
	ConnClosed                      // 连接关闭
)

var (
	ErrConnectionHanged  = errors.New("connection is hanged")
	ErrConnectionClosed  = errors.New("connection is closed")
	ErrInvalidMsgType    = errors.New("invalid message type")
	ErrTooManyConnection = errors.New("too many connection")
)

type ConnState int32

type Conn interface {
	// Cid 获取连接ID
	Cid() int64
	// Uid 获取用户ID
	Uid() int64
	// Bind 绑定用户ID
	Bind(uid int64)
	// Send 发送消息（同步）
	Send(msg []byte) error
	// Push 发送消息（异步）
	Push(msg []byte) error
	// State 获取连接状态
	State() ConnState
	// Close 关闭连接
	Close() error
	// LocalIP 获取本地IP
	LocalIP() string
	// LocalAddr 获取本地地址
	LocalAddr() net.Addr
	// RemoteIP 获取远端IP
	RemoteIP() string
	// RemoteAddr 获取远端地址
	RemoteAddr() net.Addr
}
