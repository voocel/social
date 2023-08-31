package router

const (
	Heartbeat int32 = iota + 1
	Connect
	Disconnect
	Message
	GroupMessage
)
