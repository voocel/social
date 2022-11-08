package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"strconv"
	"time"
)

type Registry struct {
	cli                            *api.Client
	Address                        string
	Token                          string
	Service                        string
	serviceID                      string
	Tag                            []string
	Port                           int
	BalanceFactor                  int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func NewRegistry(serviceName, address, token string) *Registry {
	return &Registry{
		Address:                        "127.0.0.1:8500",
		Service:                        serviceName,
		Token:                          token,
		Tag:                            []string{},
		Port:                           3000,
		BalanceFactor:                  100,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

func (r *Registry) Name() string {
	return "consul"
}

func (r *Registry) Register() error {
	config := api.DefaultConfig()
	config.Address = r.Address
	config.Token = r.Token
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	r.cli = client

	IP := localIP()
	r.serviceID = fmt.Sprintf("%v-%v-%v", r.Service, IP, r.Port)
	reg := &api.AgentServiceRegistration{
		ID:      r.serviceID, // 服务节点的名称
		Name:    r.Service,   // 服务名称
		Tags:    r.Tag,       // tag，可以为空
		Port:    r.Port,      // 服务端口
		Address: IP,          // 服务 IP
		Meta: map[string]string{
			"balanceFactor": strconv.Itoa(r.BalanceFactor),
		},
		Check: &api.AgentServiceCheck{ // 健康检查
			//TCP:                            r.Address,
			//Timeout:                        "3s",
			Interval:                       r.Interval.String(),                            // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service), // grpc支持，执行健康检查的地址，service会传到Health.Check函数中
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),      // 注销时间，相当于过期时间
		},
	}

	if err := r.cli.Agent().ServiceRegister(reg); err != nil {
		return err
	}

	return nil
}

func (r *Registry) Unregister() error {
	return r.cli.Agent().ServiceDeregister(r.serviceID)
}

func (r *Registry) QueryServices() string {
	return r.Address
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}
