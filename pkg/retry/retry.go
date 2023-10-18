package retry

import (
	"context"
	"go.uber.org/zap"
	"social/pkg/log"
	"time"
)

// Do will run function with retry mechanism.
// fn is the func to run.
// Option can control the retry times and timeout.
func Do(ctx context.Context, fn func() error, opts ...Option) error {

	c := newDefaultConfig()

	for _, opt := range opts {
		opt(c)
	}
	var el ErrorList

	for i := uint(0); i < c.attempts; i++ {
		if err := fn(); err != nil {
			if i%10 == 0 {
				log.Debug("retry func failed", zap.Uint("retry time", i), zap.Error(err))
			}

			el = append(el, err)

			if ok := IsUnRecoverable(err); ok {
				return el
			}

			select {
			case <-time.After(c.sleep):
			case <-ctx.Done():
				el = append(el, ctx.Err())
				return el
			}

			c.sleep *= 2
			if c.sleep > c.maxSleepTime {
				c.sleep = c.maxSleepTime
			}
		} else {
			return nil
		}
	}
	return el
}

type unrecoverableError struct {
	error
}

// Unrecoverable method wrap an error to unrecoverableError. This will make retry
// quick return.
func Unrecoverable(err error) error {
	return unrecoverableError{err}
}

// IsUnRecoverable is used to judge whether the error is wrapped by unrecoverableError.
func IsUnRecoverable(err error) bool {
	_, isUnrecoverable := err.(unrecoverableError)
	return isUnrecoverable
}
