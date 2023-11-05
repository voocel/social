package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var Conf = new(config)

type config struct {
	Mode            string
	Name            string
	LogLevel        string `mapstructure:"log_level"`
	LogPath         string `mapstructure:"log_path"`
	LogLevelAddr    string `mapstructure:"log_level_addr"`
	LogLevelPattern string `mapstructure:"log_level_pattern"`

	App       AppConfig
	Transport Transport
	Http      HttpConfig
	IM        IMConfig
	Group     GroupConfig
	Gateway   GateConfig
	RPC       RPCConfig
	Etcd      EtcdConfig
	Redis     RedisConfig
	Mysql     MysqlConfig
	Postgres  PostgresConfig
}

type Transport struct {
	DiscoveryNode []Discovery `mapstructure:"discovery_node"`
	DiscoveryGate string      `mapstructure:"discovery_gate"`
	Grpc          struct {
		Addr        string
		ServiceName string `mapstructure:"service_name"`
	}
}

type AppConfig struct {
	Domain         string
	StaticRootPath string `mapstructure:"static_root_path"`
}

type HttpConfig struct {
	Addr string
}

type IMConfig struct {
}

type GroupConfig struct {
}

type GateConfig struct {
	Name string
	Addr string
}

type RPCConfig struct {
}

type EtcdConfig struct {
	Addr string
}

type RedisConfig struct {
	Addr        string
	Password    string
	Db          int
	PoolSize    int `mapstructure:"pool_size"`
	MinIdleConn int `mapstructure:"min_idle_conn"`
}

type MysqlConfig struct {
	Host            string
	Port            int
	Dbname          string
	Username        string
	Password        string
	MaximumPoolSize int `mapstructure:"maximum_pool_size"`
	MaximumIdleSize int `mapstructure:"maximum_idle_size"`
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Sslmode  string
}

type Discovery struct {
	Name    string
	Routers []int32
}

func LoadConfig(paths ...string) {
	if len(paths) == 0 {
		viper.AddConfigPath(".")
		viper.AddConfigPath("config")
		viper.AddConfigPath("../config")
	} else {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("social")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("mode", "debug")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_path", "log")
	viper.SetDefault("log_level_pattern", "/log/level")
	viper.SetDefault("atomic_level_addr", "4240")

	viper.SetDefault("http.addr", ":8090")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("read config error: %v", err)
	}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Panicf("unmarshal config err: %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config change: %s, %s, %s\n", e.Op.String(), e.Name, e.String())
		if err := viper.Unmarshal(Conf); err != nil {
			log.Printf("config change unmarshal err: %v", err)
		}
	})
	log.Println("load config successfully")
}
