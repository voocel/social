package etcd

import (
	"fmt"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	prefix    string
	discovery *Discovery
}

func NewResolver(endpoints []string, prefix string) *Resolver {
	return &Resolver{
		prefix:    prefix,
		discovery: NewDiscovery(endpoints),
	}
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	key := "/" + target.URL.Scheme + "/" + target.URL.Path + "/"
	fmt.Println("resolver key: ", key)
	err := r.discovery.Watch(r.prefix)
	if err != nil {
		return nil, err
	}
	r.discovery.cc = cc
	return r, nil
}

func (r *Resolver) Scheme() string {
	return "etcd"
}

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *Resolver) Close()                                {}
