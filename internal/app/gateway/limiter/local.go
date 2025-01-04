package limiter

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type LocalLimiter struct {
	limiters sync.Map
	qps      int
	burst    int
}

func NewLocalLimiter(qps, burst int) *LocalLimiter {
	return &LocalLimiter{
		qps:   qps,
		burst: burst,
	}
}

func (l *LocalLimiter) Allow(key string) bool {
	return l.getLimiter(key).Allow()
}

func (l *LocalLimiter) AllowN(key string, n int) bool {
	return l.getLimiter(key).AllowN(time.Now(), n)
}

func (l *LocalLimiter) getLimiter(key string) *rate.Limiter {
	limiter, ok := l.limiters.Load(key)
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(l.qps), l.burst)
		l.limiters.Store(key, limiter)
	}
	return limiter.(*rate.Limiter)
}
