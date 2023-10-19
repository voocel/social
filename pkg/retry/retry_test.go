package retry

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	ctx := context.Background()

	n := 0
	testFn := func() error {
		if n < 3 {
			n++
			return fmt.Errorf("some error")
		}
		return nil
	}

	err := Do(ctx, testFn)
	assert.Nil(t, err)
}

func TestAttempts(t *testing.T) {
	ctx := context.Background()

	testFn := func() error {
		return fmt.Errorf("some error")
	}

	err := Do(ctx, testFn, Attempts(1))
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestMaxSleepTime(t *testing.T) {
	ctx := context.Background()

	testFn := func() error {
		return fmt.Errorf("some error")
	}

	err := Do(ctx, testFn, Attempts(3), MaxSleepTime(200*time.Millisecond))
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestSleep(t *testing.T) {
	ctx := context.Background()

	testFn := func() error {
		return fmt.Errorf("some error")
	}

	err := Do(ctx, testFn, Attempts(3), Sleep(500*time.Millisecond))
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestAllError(t *testing.T) {
	ctx := context.Background()

	testFn := func() error {
		return errors.New("some error")
	}

	err := Do(ctx, testFn, Attempts(3))
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestUnRecoveryError(t *testing.T) {
	attempts := 0
	ctx := context.Background()

	testFn := func() error {
		attempts++
		return Unrecoverable(fmt.Errorf("some error"))
	}

	err := Do(ctx, testFn, Attempts(3))
	assert.NotNil(t, err)
	assert.Equal(t, attempts, 1)
}

func TestContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	testFn := func() error {
		return fmt.Errorf("some error")
	}

	err := Do(ctx, testFn)
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	testFn := func() error {
		return fmt.Errorf("some error")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := Do(ctx, testFn)
	assert.NotNil(t, err)
	fmt.Println(err)
}
