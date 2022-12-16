package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

var (
	ErrNotFoundClient = errors.New("not found grpc client service")
	ErrConnShutdown   = errors.New("grpc connection has closed")

	defaultPoolSize                    = 10
	defaultDialTimeout                 = 10 * time.Second
	defaultKeepAlive                   = 30 * time.Second
	defaultKeepAliveTimeout            = 10 * time.Second
	defaultBackoffMaxDelay             = 3 * time.Second
	defaultMaxSendMsgSize              = 4 << 20
	defaultMaxMaxRecvMsgSize           = 4 << 20
	defaultInitialWindowSize     int32 = 4 << 20
	defaultInitialConnWindowSize int32 = 4 << 20
)

type Option struct {
	PoolSize         int
	DialTimeOut      time.Duration
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
	DialOptions      []grpc.DialOption
}

type Pool struct {
	endpoint string
	next     int64
	cap      int64

	option *Option
	conns  []*grpc.ClientConn
	sync.Mutex
}

func (p *Pool) getConn() (*grpc.ClientConn, error) {
	var (
		idx  int64
		next int64
		err  error
	)

	next = atomic.AddInt64(&p.next, 1)
	idx = next % p.cap
	conn := p.conns[idx]
	if conn != nil && p.checkState(conn) == nil {
		return conn, nil
	}

	if conn != nil {
		conn.Close()
	}

	p.Lock()
	defer p.Unlock()

	// double check to prevent initialization
	conn = p.conns[idx]
	if conn != nil && p.checkState(conn) == nil {
		return conn, nil
	}

	conn, err = p.connect()
	if err != nil {
		return nil, err
	}

	p.conns[idx] = conn
	return conn, nil
}

func (p *Pool) checkState(conn *grpc.ClientConn) error {
	state := conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return ErrConnShutdown
	}

	return nil
}

// refer to https://github.com/grpc/grpc-proto/blob/master/grpc/service_config/service_config.proto
//retryPolicy := fmt.Sprintf(`{
//		"methodConfig": [{
//		  "name": [{"service": "%s"}],
//		  "retryPolicy": {
//			  "MaxAttempts": %d,
//			  "InitialBackoff": "%fs",
//			  "MaxBackoff": "%fs",
//			  "BackoffMultiplier": %f,
//			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
//		  }
//		}]}`, RetryServiceName, MaxAttempts, InitialBackoff, MaxBackoff, BackoffMultiplier)
func (p *Pool) defaultDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		//grpc.WithBlock(),
		//grpc.WithConnectParams(grpc.ConnectParams{
		//	Backoff:           backoff.Config{MaxDelay: 8*time.Second},
		//	MinConnectTimeout: 0,
		//}),
		//grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{
			"LoadBalancingPolicy": "%s",
			"MethodConfig": [
				{
					"Name": [{"Service": "helloworld.Greeter"}], 
					"RetryPolicy": {
						"MaxAttempts":2, "InitialBackoff": "0.1s", 
						"MaxBackoff": "1s", "BackoffMultiplier": 2.0, 
						"RetryableStatusCodes": ["UNAVAILABLE", "CANCELLED"]
					}
				}
			],
			"HealthCheckConfig": {"ServiceName": "helloworld.Greeter"}
		}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithInitialWindowSize(defaultInitialWindowSize),
		grpc.WithInitialConnWindowSize(defaultInitialConnWindowSize),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(defaultMaxSendMsgSize),
			grpc.MaxCallRecvMsgSize(defaultMaxMaxRecvMsgSize),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                p.option.KeepAlive,
			Timeout:             p.option.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   3 * time.Second,
			},
			MinConnectTimeout: defaultDialTimeout,
		}),
	}
}

func (p *Pool) connect() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), p.option.DialTimeOut)
	defer cancel()
	opts := p.option.DialOptions
	if opts == nil {
		opts = p.defaultDialOptions()
	}
	conn, err := grpc.DialContext(ctx, p.endpoint, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (p *Pool) Close() {
	p.Lock()
	defer p.Unlock()

	for _, conn := range p.conns {
		if conn == nil {
			continue
		}
		conn.Close()
	}
}

func newClientPoolWithOption(endpoint string, option *Option) *Pool {
	if (option.PoolSize) <= 0 {
		option.PoolSize = defaultPoolSize
	}

	if option.DialTimeOut <= 0 {
		option.DialTimeOut = defaultDialTimeout
	}

	if option.KeepAlive <= 0 {
		option.KeepAlive = defaultKeepAlive
	}

	if option.KeepAliveTimeout <= 0 {
		option.KeepAliveTimeout = defaultKeepAliveTimeout
	}

	return &Pool{
		endpoint: endpoint,
		option:   option,
		cap:      int64(option.PoolSize),
		conns:    make([]*grpc.ClientConn, option.PoolSize),
	}
}

type ServiceClientPool struct {
	useTLS   bool
	option   *Option
	services map[string][]string
	clients  map[string]*Pool
}

func NewServiceClientPool(option *Option) *ServiceClientPool {
	return &ServiceClientPool{
		option:   option,
		services: make(map[string][]string),
	}
}

func (sc *ServiceClientPool) Start() {
	var clients = make(map[string]*Pool, len(sc.services))
	if !sc.useTLS {
		sc.addDialOption(grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	for endpoint, serviceName := range sc.services {
		clientPool := newClientPoolWithOption(endpoint, sc.option)
		for _, srv := range serviceName {
			clients[srv] = clientPool
		}
	}

	sc.clients = clients
	scp = sc
}

func (sc *ServiceClientPool) addDialOption(opt grpc.DialOption) {
	sc.option.DialOptions = append(sc.option.DialOptions, opt)
}

// SetUnaryInterceptors Call once, multiple calls, only the last one takes effect
func (sc *ServiceClientPool) SetUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) {
	chain := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...))
	sc.addDialOption(chain)
}

func (sc *ServiceClientPool) SetStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) {
	chain := grpc.WithChainStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...))
	sc.addDialOption(chain)
}

// SetTLS mutual TLS
func (sc *ServiceClientPool) SetTLS(clientKeyPath, clientPemPath, caPemPath, commonName string) {
	opt, err := setCert(clientKeyPath, clientPemPath, caPemPath, commonName)
	if err != nil {
		panic(err)
	}
	sc.useTLS = true
	sc.addDialOption(opt)
}

func (sc *ServiceClientPool) SetServices(endpoint string, services ...string) {
	if len(services) == 0 {
		return
	}
	sc.services[endpoint] = append(sc.services[endpoint], services...)
}

func (sc *ServiceClientPool) GetClientWithFullMethod(fullMethod string) (*grpc.ClientConn, error) {
	sn := sc.spiltFullMethod(fullMethod)
	return sc.GetClient(sn)
}

func (sc *ServiceClientPool) GetClient(serviceName string) (*grpc.ClientConn, error) {
	cc, ok := sc.clients[serviceName]
	if !ok {
		return nil, ErrNotFoundClient
	}

	return cc.getConn()
}

func (sc *ServiceClientPool) Close(serviceName string) {
	cc, ok := sc.clients[serviceName]
	if !ok {
		return
	}

	cc.Close()
}

func (sc *ServiceClientPool) CloseAll() {
	for _, client := range sc.clients {
		client.Close()
	}
}

func (sc *ServiceClientPool) spiltFullMethod(fullMethod string) string {
	var arr []string
	arr = strings.Split(fullMethod, "/")
	if len(arr) != 3 {
		return ""
	}

	return arr[1]
}

func (sc *ServiceClientPool) Invoke(
	ctx context.Context,
	fullMethod string,
	headers map[string]string,
	args interface{},
	reply interface{},
	opts ...grpc.CallOption,
) error {
	var md metadata.MD
	serviceName := sc.spiltFullMethod(fullMethod)
	conn, err := sc.GetClient(serviceName)
	if err != nil {
		return err
	}

	md, exist := metadata.FromOutgoingContext(ctx)
	if exist {
		md = md.Copy()
	} else {
		md = metadata.MD{}
	}

	for k, v := range headers {
		md.Set(k, v)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	return conn.Invoke(ctx, fullMethod, args, reply, opts...)
}

var scp *ServiceClientPool

func NewDefaultPool() *ServiceClientPool {
	co := Option{
		PoolSize:         defaultPoolSize,
		DialTimeOut:      defaultDialTimeout,
		KeepAlive:        defaultKeepAlive,
		KeepAliveTimeout: defaultKeepAliveTimeout,
	}
	return NewServiceClientPool(&co)
}

func GetPool() *ServiceClientPool {
	return scp
}
