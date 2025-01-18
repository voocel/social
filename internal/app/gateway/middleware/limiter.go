package middleware

import (
	"social/internal/app/gateway/limiter"
	"social/pkg/network"
)

type RateLimiterMiddleware struct {
	config        *limiter.Config
	globalLimiter limiter.RateLimiter
	ipLimiter     limiter.RateLimiter
	userLimiter   limiter.RateLimiter
	routeLimiters map[string]limiter.RateLimiter
}

func NewRateLimiterMiddleware(conf *limiter.Config) *RateLimiterMiddleware {
	m := &RateLimiterMiddleware{
		config: conf,
		globalLimiter: limiter.NewRedisLimiter(
			conf.Rules.Global.QPS,
			conf.Rules.Global.Burst,
		),
		ipLimiter: limiter.NewRedisLimiter(
			conf.Rules.IP.QPS,
			conf.Rules.IP.Burst,
		),
		userLimiter: limiter.NewRedisLimiter(
			conf.Rules.User.QPS,
			conf.Rules.User.Burst,
		),
		routeLimiters: make(map[string]limiter.RateLimiter),
	}

	// 初始化路由限流器
	for route, rule := range conf.Rules.API.Routes {
		m.routeLimiters[route] = limiter.NewRedisLimiter(rule.QPS, rule.Burst)
	}

	return m
}

func (m *RateLimiterMiddleware) Handle(next network.HandlerFunc) network.HandlerFunc {
	return func(conn network.Conn) {
		if !m.config.Enable {
			next(conn)
			return
		}

		// IP限流
		clientIP := conn.RemoteAddr().String()
		if !m.ipLimiter.Allow(clientIP) {
			conn.Close()
			return
		}

		// 用户限流
		if uid := conn.Uid(); uid > 0 {
			if !m.userLimiter.Allow(fmt.Sprintf("user:%d", uid)) {
				conn.Close()
				return
			}
		}

		next(conn)
	}
}
