package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"regexp"
	"sync"
)

const (
	defaultPort = "8000"
)

var regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")

type Builder struct {
	name    string
	address string
	token   string
}

func NewResolver(serviceName, address, token string) string {
	builder := &Builder{
		name:    serviceName,
		address: address,
		token:   token,
	}
	resolver.Register(builder)

	return "consul:///" + serviceName
}

type Resolver struct {
	address              string
	service              string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	lastIndex            uint64
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	host, port, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	if err != nil {
		return nil, err
	}

	cr := &Resolver{
		address:              fmt.Sprintf("%s%s", host, port),
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	cr.wg.Add(1)
	go cr.watcher()
	return cr, nil
}

func (b *Builder) Scheme() string {
	return "consul"
}

func (r *Resolver) ResolveNow(opt resolver.ResolveNowOptions) {}

func (r *Resolver) Close() {}

func (r *Resolver) watcher() {
	fmt.Printf("calling consul watcher\n")
	config := api.DefaultConfig()
	config.Address = r.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}

	for {
		services, metainfo, err := client.Health().Service(r.name, r.name, true, &api.QueryOptions{WaitIndex: r.lastIndex})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		}

		r.lastIndex = metainfo.LastIndex
		var newAddrs []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			newAddrs = append(newAddrs, resolver.Address{Addr: addr})
		}
		fmt.Printf("add service addrs\n")
		fmt.Printf("newAddrs: %v\n", newAddrs)
		r.cc.UpdateState(resolver.State{Addresses: newAddrs})
	}

}

func parseTarget(target string) (host, port, name string, err error) {
	if target == "" {
		err = errors.New("consul resolver: missing address")
		return
	}

	if !regexConsul.MatchString(target) {
		err = errors.New("consul resolver: invalid uri")
		return
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	if port == "" {
		port = defaultPort
	}
	return host, port, name, nil
}
