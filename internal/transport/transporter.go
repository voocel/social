package transport

type Transporter interface {
	Options() *Options
	// NewGateServer 新建网关服务器
	NewGateServer(provider GateProvider) Server
	// NewNodeServer 新建节点服务器
	NewNodeServer(provider NodeProvider) Server
	// NewGateClient 新建网关客户端
	NewGateClient(addr string) (GateClient, error)
	// NewNodeClient 新建节点客户端
	NewNodeClient(addr string) (NodeClient, error)
}

type Message struct {
	Seq    int32  // 序列号
	Route  int32  // 路由
	Buffer []byte // 消息内容
}
