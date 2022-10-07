package discovery

type Discovery interface {
	// QueryServices 向注册中心查询所有服务
	QueryServices() ([]*ServiceInfo, error)

	// Register 注册服务
	Register() error

	// UnRegister 取消服务
	UnRegister() error

	// GetService 获取指定服务
	GetService(name string) []ServiceInfo
}

type Watcher interface {
	// Next 返回服务实例列表
	Next() ([]*ServiceInfo, error)
	// Stop 停止监听
	Stop() error
}

type ServiceInfo struct {
	ServiceName string
	Addr        string
	Meta        map[string]string
}
