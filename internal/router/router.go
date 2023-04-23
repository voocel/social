package router

import "sync"

type Route struct {
	id              int32    // 路由ID
	stateful        bool     // 是否有状态
	endpoints       sync.Map // 服务端口
	balanceStrategy string   // 负载均衡策略
}

type Router struct {
	rw            sync.RWMutex
	routes        map[int32]*Route  // 节点路由表
	gateEndpoints map[string]string // 网关服务端口
	nodeEndpoints map[string]string // 节点服务端口
}

func NewRouter() *Router {
	return &Router{
		routes:        make(map[int32]*Route),
		gateEndpoints: make(map[string]string),
		nodeEndpoints: make(map[string]string),
	}
}
