package transport

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
