package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"social/config"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		DB:           config.Conf.Redis.Db,
		Addr:         config.Conf.Redis.Addr,
		Password:     config.Conf.Redis.Password,
		PoolSize:     config.Conf.Redis.PoolSize,
		MinIdleConns: config.Conf.Redis.MinIdleConn,
	})
	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
}

func GetRedisClient() *redis.Client {
	if nil == redisClient {
		panic("Please initialize the Redis client first!")
	}
	return redisClient
}

func CloseRedis() {
	if nil != redisClient {
		_ = redisClient.Close()
	}
}
