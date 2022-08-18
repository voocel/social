package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	LogLevel string
	LogPath  string

	Http  HttpConfig
	IM    IMConfig
	Gate  GateConfig
	RPC   RPCConfig
	Redis RedisConfig
}

type HttpConfig struct {
	Addr string
}

type IMConfig struct {
}

type GateConfig struct {
}

type RPCConfig struct {
}

type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	Db          int
	PoolSize    int
	MinIdleConn int
}

func LoadConfig(path ...string) {
	if len(path) == 0 {
		viper.AddConfigPath("config")
	} else {
		viper.AddConfigPath(path[0])
	}

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("social")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_path", "info")
	viper.SetDefault("http.addr", ":8090")
	viper.SetDefault("gate", nil)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("load config error: %v", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("file change: %s, %s\n", e.Name, viper.GetString("websocket.port"))
	})
	log.Println("load config successfully")
}

func GetConfig() *Config {
	return &Config{
		LogLevel: viper.GetString("log_level"),
		LogPath:  viper.GetString("log_path"),
		Http:     GetHttpConfig(),
		IM:       GetIMConfig(),
		Redis:    GetRedisConfig(),
	}
}

func GetHttpConfig() HttpConfig {
	return HttpConfig{
		Addr: viper.GetString("http.addr"),
	}
}

func GetIMConfig() IMConfig {
	return IMConfig{}
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
