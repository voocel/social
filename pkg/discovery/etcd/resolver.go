package etcd

import (
	"google.golang.org/grpc/resolver"
	"log"
)

type Resolver struct {
	prefix   string
	registry *Registry
}

func NewResolver(endpoints []string, prefix string) (*Resolver, error) {
	r, err := NewRegistry(endpoints)
	return &Resolver{
		prefix:   prefix,
		registry: r,
	}, err
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	key := target.URL.Scheme + target.URL.Path + "/"
	log.Println("resolver key: ", key)
	err := r.registry.watch(key)
	if err != nil {
		return nil, err
	}
	r.registry.cc = cc
	return r, nil
}

func (r *Resolver) Scheme() string {
	return "etcd"
}

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *Resolver) Close()                                {}
