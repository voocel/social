package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var Conf = new(Config)

type Config struct {
	Mode            string
	LogLevel        string `mapstructure:"log_level"`
	LogPath         string `mapstructure:"log_path"`
	AtomicLevelAddr string `mapstructure:"atomic_level_addr"`

	Transport Transport
	Http      HttpConfig
	IM        IMConfig
	Group     GroupConfig
	Gate      GateConfig
	RPC       RPCConfig
	Redis     RedisConfig
	Mysql     MysqlConfig
}

type Transport struct {
	DiscoveryNode []Discovery `mapstructure:"discovery_node"`
	DiscoveryGate string      `mapstructure:"discovery_gate"`
	Grpc          struct {
		Addr        string
		ServiceName string `mapstructure:"service_name"`
	}
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

type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	Db          int
	PoolSize    int
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
	viper.SetDefault("atomic_level_addr", "4240")

	viper.SetDefault("http.addr", ":8090")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("load config error: %v", err)
		return
	}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Panic(err)
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

func NewConfig() *Config {
	return &Config{
		LogLevel: viper.GetString("log_level"),
		LogPath:  viper.GetString("log_path"),
		Mode:     viper.GetString("mode"),
		Http:     GetHttpConfig(),
		IM:       GetIMConfig(),
		Group:    GetGroupConfig(),
		Redis:    GetRedisConfig(),
		Mysql:    GetMysqlConfig(),
	}
}

func (c *Config) Get(key string) string {
	return viper.GetString(key)
}

func GetHttpConfig() HttpConfig {
	return HttpConfig{
		Addr: viper.GetString("http.addr"),
	}
}

func GetIMConfig() IMConfig {
	return IMConfig{}
}

func GetGroupConfig() GroupConfig {
	return GroupConfig{}
}

func GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:        viper.GetString("redis.host"),
		Port:        viper.GetInt("redis.port"),
		Password:    viper.GetString("redis.password"),
		Db:          viper.GetInt("redis.db"),
		PoolSize:    viper.GetInt("redis.pool_size"),
		MinIdleConn: viper.GetInt("redis.min_idle_conn"),
	}
}

func GetMysqlConfig() MysqlConfig {
	return MysqlConfig{
		Host:            viper.GetString("mysql.host"),
		Port:            viper.GetInt("mysql.port"),
		Dbname:          viper.GetString("mysql.dbname"),
		Username:        viper.GetString("mysql.username"),
		Password:        viper.GetString("mysql.password"),
		MaximumPoolSize: viper.GetInt("mysql.maximum_pool_size"),
		MaximumIdleSize: viper.GetInt("mysql.maximum_idle_size"),
	}
}
