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
	Mode     string

	Http  HttpConfig
	IM    IMConfig
	Gate  GateConfig
	RPC   RPCConfig
	Redis RedisConfig
	Mysql MysqlConfig
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

type MysqlConfig struct {
	Host            string
	Port            int
	Dbname          string
	Username        string
	Password        string
	MaximumPoolSize int
	MaximumIdleSize int
}

func LoadConfig(paths ...string) {
	if len(paths) == 0 {
		viper.AddConfigPath(".")
		viper.AddConfigPath("config")
		viper.AddConfigPath("../config")
		viper.AddConfigPath("../../config")
	} else {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("social")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_path", "info")
	viper.SetDefault("mode", "debug")
	viper.SetDefault("http.addr", ":8090")
	viper.SetDefault("gate", nil)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("load config error: %v", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("file change: %s, %s, %s\n", e.Op.String(), e.Name, e.String())
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
