package nacos

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"social/pkg/log"
)

// NacosRegisterPlugin implements consul registry.
type NacosRegisterPlugin struct {
	// service address, for example, tcp@127.0.0.1:8972, quic@127.0.0.1:1234
	ServiceAddress string
	// nacos client config
	ClientConfig constant.ClientConfig
	// nacos server config
	ServerConfig []constant.ServerConfig
	Cluster      string
	Group        string
	Weight       float64

	// Registered services
	Services []string

	namingClient naming_client.INamingClient
}

// Start starts to connect consul cluster
func (p *NacosRegisterPlugin) Start() error {
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"clientConfig":  p.ClientConfig,
		"serverConfigs": p.ServerConfig,
	})
	if err != nil {
		return err
	}

	p.namingClient = namingClient

	return nil
}

// Stop unregister all services.
func (p *NacosRegisterPlugin) Stop() error {
	_, ip, port, _ := ParseAddress(p.ServiceAddress)

	for _, name := range p.Services {
		inst := vo.DeregisterInstanceParam{
			Ip:          ip,
			Ephemeral:   true,
			Port:        uint64(port),
			ServiceName: name,
			Cluster:     p.Cluster,
			GroupName:   p.Group,
		}
		_, err := p.namingClient.DeregisterInstance(inst)
		if err != nil {
			log.Errorf("faield to deregister %s: %v", name, err)
		}
	}

	return nil
}

// Register handles registering event.
// this service is registered at BASE/serviceName/thisIpAddress node
func (p *NacosRegisterPlugin) Register(name string, rcvr interface{}, metadata string) (err error) {
	if strings.TrimSpace(name) == "" {
		return errors.New("Register service `name` can't be empty")
	}

	network, ip, port, err := ParseAddress(p.ServiceAddress)
	if err != nil {
		log.Errorf("failed to parse rpcx addr in Register: %v", err)
		return err
	}

	meta := ConvertMeta2Map(metadata)
	meta["network"] = network

	inst := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: name,
		Metadata:    meta,
		ClusterName: p.Cluster,
		GroupName:   p.Group,
		Weight:      p.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}

	_, err = p.namingClient.RegisterInstance(inst)
	if err != nil {
		log.Errorf("failed to register %s: %v", name, err)
		return err
	}

	p.Services = append(p.Services, name)

	return
}

func (p *NacosRegisterPlugin) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return p.Register(serviceName, fn, metadata)
}

func (p *NacosRegisterPlugin) Unregister(name string) (err error) {
	if len(p.Services) == 0 {
		return nil
	}

	if strings.TrimSpace(name) == "" {
		return errors.New("unregister service `name` can't be empty")
	}

	_, ip, port, err := ParseAddress(p.ServiceAddress)
	if err != nil {
		log.Errorf("wrong address %s: %v", p.ServiceAddress, err)
		return err
	}

	inst := vo.DeregisterInstanceParam{
		Ip:          ip,
		Ephemeral:   true,
		Port:        uint64(port),
		ServiceName: name,
		Cluster:     p.Cluster,
		GroupName:   p.Group,
	}
	_, err = p.namingClient.DeregisterInstance(inst)
	if err != nil {
		log.Errorf("failed to deregister %s: %v", name, err)
		return err
	}

	services := make([]string, 0, len(p.Services)-1)
	for _, s := range p.Services {
		if s != name {
			services = append(services, s)
		}
	}
	p.Services = services

	return nil
}

// ParseAddress parses rpcx address such as tcp@127.0.0.1:8972  quic@192.168.1.1:9981
func ParseAddress(addr string) (network string, ip string, port int, err error) {
	ati := strings.Index(addr, "@")
	if ati <= 0 {
		return "", "", 0, fmt.Errorf("invalid rpcx address: %s", addr)
	}

	network = addr[:ati]
	addr = addr[ati+1:]

	var portstr string
	ip, portstr, err = net.SplitHostPort(addr)
	if err != nil {
		return "", "", 0, err
	}

	port, err = strconv.Atoi(portstr)
	return network, ip, port, err
}

func ConvertMeta2Map(meta string) map[string]string {
	var rt = make(map[string]string)

	if meta == "" {
		return rt
	}

	v, err := url.ParseQuery(meta)
	if err != nil {
		return rt
	}

	for key := range v {
		rt[key] = v.Get(key)
	}
	return rt
}
