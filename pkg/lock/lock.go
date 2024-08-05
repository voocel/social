package lock

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

const (
	ttl              = time.Second * 30
	resetTTLInterval = ttl / 3
	tryLockInterval  = time.Second
	unlockScript     = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`
)

var (
	// ErrLockFailed Lock failed
	ErrLockFailed = errors.New("lock failed")
	// ErrTimeout Lock timeout
	ErrTimeout = errors.New("timeout")
)

type Locker struct {
	client          *redis.Client
	script          *redis.Script
	ttl             time.Duration
	tryLockInterval time.Duration
}

func NewLocker(client *redis.Client, ttl, tryLockInterval time.Duration) *Locker {
	return &Locker{
		client:          client,
		script:          redis.NewScript(unlockScript),
		ttl:             ttl,
		tryLockInterval: tryLockInterval,
	}
}

func (l *Locker) GetLock(resource string) *Lock {
	return &Lock{
		client:          l.client,
		script:          l.script,
		resource:        resource,
		randomValue:     uuid.New().String(),
		watchDog:        make(chan struct{}),
		ttl:             l.ttl,
		tryLockInterval: l.tryLockInterval,
	}
}

// Lock Non reusable
type Lock struct {
	client          *redis.Client
	script          *redis.Script
	resource        string
	randomValue     string
	watchDog        chan struct{}
	ttl             time.Duration
	tryLockInterval time.Duration
}

func (l *Lock) Lock(ctx context.Context) error {
	// Try adding a lock
	err := l.TryLock(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrLockFailed) {
		return err
	}
	// Lock failed, keep trying
	ticker := time.NewTicker(l.tryLockInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// Timeout
			return ErrTimeout
		case <-ticker.C:
			// Try locking again
			err := l.TryLock(ctx)
			if err == nil {
				return nil
			}
			if !errors.Is(err, ErrLockFailed) {
				return err
			}
		}
	}
}

func (l *Lock) TryLock(ctx context.Context) error {
	success, err := l.client.SetNX(ctx, l.resource, l.randomValue, l.ttl).Result()
	if err != nil {
		return err
	}

	if !success {
		return ErrLockFailed
	}
	// Lock successful, start watchdog
	go l.startWatchDog()
	return nil
}

func (l *Lock) startWatchDog() {
	ticker := time.NewTicker(l.ttl / 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Extend the expiration time of the lock
			ctx, cancel := context.WithTimeout(context.Background(), l.ttl/3*2)
			ok, err := l.client.Expire(ctx, l.resource, l.ttl).Result()
			cancel()
			// If the exception or lock no longer exists, it will not be renewed
			if err != nil || !ok {
				return
			}
		case <-l.watchDog:
			// has unlocked
			return
		}
	}
}

func (l *Lock) Unlock(ctx context.Context) error {
	err := l.script.Run(ctx, l.client, []string{l.resource}, l.randomValue).Err()
	close(l.watchDog)
	return err
}
