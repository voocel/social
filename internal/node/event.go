package node

// Event 事件
type Event int

const (
	Reconnect  Event = iota + 1 // 断线重连
	Disconnect                  // 断开连接
)

type EventHandler func(gid string, uid int64)

type eventEntity struct {
	event Event
	gid   string
	uid   int64
}
