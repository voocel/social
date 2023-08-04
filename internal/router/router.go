package router

import (
	"errors"
	"social/pkg/discovery"
	"sync"
)

var (
	ErrNotFoundRoute    = errors.New("not found route")
	ErrNotFoundEndpoint = errors.New("not found endpoint")
)

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

func (r *Router) AddService(service *discovery.Node) {

}

func (r *Router) FindGateEndpoint(gid string) (string, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	ep, ok := r.gateEndpoints[gid]
	if !ok {
		return "", ErrNotFoundEndpoint
	}

	return ep, nil
}
