package route

const (
	Heartbeat int32 = iota + 1
	Auth
	Connect
	Disconnect
	Message
	GroupMessage
	System
)

var RouteMap = map[int32]string{
	Heartbeat:    "heartbeat",
	Auth:         "auth",
	Connect:      "connect",
	Disconnect:   "disconnect",
	Message:      "message",
	GroupMessage: "group_message",
	System:       "system",
}
