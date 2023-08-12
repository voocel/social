package transport

type Transporter interface {
	// NewGateServer 新建网关服务器
	NewGateServer() (Server, error)
	// NewNodeServer 新建节点服务器
	NewNodeServer() (Server, error)
	// NewGateClient 新建网关客户端
	NewGateClient(addr string) (GateClient, error)
	// NewNodeClient 新建节点客户端
	NewNodeClient(addr string) (NodeClient, error)
}
