package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"social/config"
)

var redisClient *Redis

type Redis struct {
	client *redis.Client
}

func Init() {
	c := redis.NewClient(&redis.Options{
		DB:           config.Conf.Redis.Db,
		Addr:         config.Conf.Redis.Addr,
		Password:     config.Conf.Redis.Password,
		PoolSize:     config.Conf.Redis.PoolSize,
		MinIdleConns: config.Conf.Redis.MinIdleConn,
	})
	_, err := c.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	redisClient.client = c
}

func GetClient() *Redis {
	if nil == redisClient {
		panic("Please initialize the Redis client first!")
	}
	return redisClient
}

func Close() {
	if nil != redisClient {
		_ = redisClient.client.Close()
	}
}
