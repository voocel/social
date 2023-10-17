package etcd

import (
	"go.uber.org/zap"
	"net/url"
	"social/config"
	"social/pkg/log"
	"strconv"
	"strings"
	"sync"
)

var pulsarOnce sync.Once
var defaultYaml = "social.yaml"

const (
	// SuggestPulsarMaxMessageSize defines the maximum size of Pulsar message.
	SuggestPulsarMaxMessageSize = 5 * 1024 * 1024
	defaultEtcdLogLevel         = "info"
	defaultEtcdLogPath          = "stdout"
	KafkaProducerConfigPrefix   = "kafka.producer."
	KafkaConsumerConfigPrefix   = "kafka.consumer."
)

// BaseTable the basics of paramtable
type BaseTable struct {
	once      sync.Once
	configDir string
	RoleName  string
	YamlFile  string
	cfg       config.Config
}

// LoadWithDefault loads an object with @key. If the object does not exist, @defaultValue will be returned.
func (b *BaseTable) LoadWithDefault(key, defaultValue string) string {
	str := b.cfg.Get(key)
	if str == "" {
		return defaultValue
	}
	return str
}

// ServiceParam is used to quickly and easily access all basic service configurations.
type ServiceParam struct {
	BaseTable
	EtcdCfg   EtcdConfig
	PulsarCfg PulsarConfig
}

func (s *ServiceParam) Init() {

}

type EtcdConfig struct {
	Base *BaseTable

	// --- ETCD ---
	Endpoints         []string
	MetaRootPath      string
	KvRootPath        string
	EtcdLogLevel      string
	EtcdLogPath       string
	EtcdUseSSL        bool
	EtcdTLSCert       string
	EtcdTLSKey        string
	EtcdTLSCACert     string
	EtcdTLSMinVersion string

	// --- Embed ETCD ---
	UseEmbedEtcd bool
	ConfigPath   string
	DataDir      string
}

type PulsarConfig struct {
	Base *BaseTable

	Address        string
	WebAddress     string
	MaxMessageSize int
}

func (p *PulsarConfig) init(base *BaseTable) {
	p.Base = base

	p.initAddress()
	p.initWebAddress()
	p.initMaxMessageSize()
}

func (p *PulsarConfig) initWebAddress() {
	if p.Address == "" {
		return
	}

	pulsarURL, err := url.ParseRequestURI(p.Address)
	if err != nil {
		p.WebAddress = ""
		log.Info("failed to parse pulsar config, assume pulsar not used", zap.Error(err))
	} else {
		webport := p.Base.LoadWithDefault("pulsar.webport", "80")
		p.WebAddress = "http://" + pulsarURL.Hostname() + ":" + webport
	}
	pulsarOnce.Do(func() {
		//cmdutils.PulsarCtlConfig.WebServiceURL = p.WebAddress
	})
}

func (p *PulsarConfig) initAddress() {
	pulsarHost := p.Base.LoadWithDefault("pulsar.address", "")
	if strings.Contains(pulsarHost, ":") {
		p.Address = pulsarHost
		return
	}

	port := p.Base.LoadWithDefault("pulsar.port", "")
	if len(pulsarHost) != 0 && len(port) != 0 {
		p.Address = "pulsar://" + pulsarHost + ":" + port
	}
}

func (p *PulsarConfig) initMaxMessageSize() {
	maxMessageSizeStr := p.Base.cfg.Get("pulsar.maxMessageSize")
	if maxMessageSizeStr == "" {
		p.MaxMessageSize = SuggestPulsarMaxMessageSize
	} else {
		maxMessageSize, err := strconv.Atoi(maxMessageSizeStr)
		if err != nil {
			p.MaxMessageSize = SuggestPulsarMaxMessageSize
		} else {
			p.MaxMessageSize = maxMessageSize
		}
	}
}
