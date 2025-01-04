package limiter

import (
	"context"
	"social/pkg/redis"
)

type RedisLimiter struct {
	redis  *redis.Redis
	qps    int
	burst  int
	prefix string
}

func NewRedisLimiter(qps, burst int) *RedisLimiter {
	return &RedisLimiter{
		redis:  redis.GetClient(),
		qps:    qps,
		burst:  burst,
		prefix: "rate_limit:",
	}
}

func (l *RedisLimiter) Allow(key string) bool {
	return l.AllowN(key, 1)
}

func (l *RedisLimiter) AllowN(key string, n int) bool {
	lua := `
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local window = tonumber(ARGV[2])
        local increment = tonumber(ARGV[3])
        local current = redis.call('INCRBY', key, increment)
        if current == increment then
            redis.call('EXPIRE', key, window)
        end
        return current <= limit
    `

	result, err := l.redis.Eval(context.Background(), lua,
		[]string{l.prefix + key},
		l.burst,
		1, // 1秒窗口
		n,
	).Result()

	if err != nil {
		return false
	}
	return result.(int64) == 1
}
