package node

type Message struct {
	Seq   int32       // 序列号
	Route int32       // 路由ID
	Data  interface{} // 消息数据，接收json、proto、[]byte
}

type DeliverArgs struct {
	NID     string   // 接收节点。存在接收节点时，消息会直接投递给接收节点；不存在接收节点时，系统定位用户所在节点，然后投递。
	UID     int64    // 用户ID
	Message *Message // 消息
}
