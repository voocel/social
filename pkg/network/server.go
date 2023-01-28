package network

type (
	StartHandler      func()
	CloseHandler      func()
	ConnectHandler    func(conn Conn)
	DisconnectHandler func(conn Conn, err error)
	ReceiveHandler    func(conn Conn, msg []byte)
)

type Server interface {
	// Start 启动服务器
	Start() error
	// Stop 关闭服务器
	Stop() error
	// Protocol 协议
	Protocol() string
	// OnStart 监听服务器启动
	OnStart(handler StartHandler)
	// OnStop 监听服务器关闭
	OnStop(handler CloseHandler)
	// OnConnect 监听连接打开
	OnConnect(handler ConnectHandler)
	// OnReceive 监听接收消息
	OnReceive(handler ReceiveHandler)
	// OnDisconnect 监听连接断开
	OnDisconnect(handler DisconnectHandler)
}
