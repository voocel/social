package lru

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLRU(t *testing.T) {
	c, err := NewLRU(1, nil)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	c, err = NewLRU(0, nil)
	assert.NotNil(t, err)
	assert.Nil(t, c)

	c, err = NewLRU(-1, nil)
	assert.NotNil(t, err)
	assert.Nil(t, c)
}

func TestLRU_Add(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testValueExtra := "test_value_extra"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"

	testKey3 := "test_key_3"
	testValue3 := "test_value_3"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)

	v, ok := c.Get(testKey1)
	assert.True(t, ok)
	assert.EqualValues(t, testValue1, v)

	v, ok = c.Get(testKey2)
	assert.True(t, ok)
	assert.EqualValues(t, testValue2, v)

	c.Add(testKey1, testValueExtra)

	k, v, ok := c.GetOldest()
	assert.True(t, ok)
	assert.EqualValues(t, testKey2, k)
	assert.EqualValues(t, testValue2, v)

	c.Add(testKey3, testValue3)
	v, ok = c.Get(testKey3)
	assert.True(t, ok)
	assert.EqualValues(t, testValue3, v)

	v, ok = c.Get(testKey2)
	assert.False(t, ok)
	assert.Nil(t, v)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_Contains(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(1, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"

	c.Add(testKey1, testValue1)
	ok := c.Contains(testKey1)
	assert.True(t, ok)

	c.Add(testKey2, testValue2)
	ok = c.Contains(testKey2)
	assert.True(t, ok)

	ok = c.Contains(testKey1)
	assert.False(t, ok)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_Get(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(1, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"

	c.Add(testKey1, testValue1)
	v, ok := c.Get(testKey1)
	assert.True(t, ok)
	assert.EqualValues(t, testValue1, v)

	c.Add(testKey2, testValue2)
	v, ok = c.Get(testKey2)
	assert.True(t, ok)
	assert.EqualValues(t, testValue2, v)

	v, ok = c.Get(testKey1)
	assert.False(t, ok)
	assert.Nil(t, v)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_GetOldest(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"
	testKey3 := "test_key_3"
	testValue3 := "test_value_3"

	k, v, ok := c.GetOldest()
	assert.False(t, ok)
	assert.Nil(t, k)
	assert.Nil(t, v)

	c.Add(testKey1, testValue1)
	k, v, ok = c.GetOldest()
	assert.True(t, ok)
	assert.EqualValues(t, testKey1, k)
	assert.EqualValues(t, testValue1, v)

	c.Add(testKey2, testValue2)
	k, v, ok = c.GetOldest()
	assert.True(t, ok)
	assert.EqualValues(t, testKey1, k)
	assert.EqualValues(t, testValue1, v)

	v, ok = c.Get(testKey1)
	assert.True(t, ok)
	assert.EqualValues(t, testValue1, v)

	k, v, ok = c.GetOldest()
	assert.True(t, ok)
	assert.EqualValues(t, testKey2, k)
	assert.EqualValues(t, testValue2, v)

	c.Add(testKey3, testValue3)
	k, v, ok = c.GetOldest()
	assert.True(t, ok)
	assert.EqualValues(t, testKey1, k)
	assert.EqualValues(t, testValue1, v)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_Keys(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"
	testKey3 := "test_key_3"
	testValue3 := "test_value_3"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)
	keys := c.Keys()
	assert.ElementsMatch(t, []string{testKey1, testKey2}, keys)

	v, ok := c.Get(testKey1)
	assert.True(t, ok)
	assert.EqualValues(t, testValue1, v)

	keys = c.Keys()
	assert.ElementsMatch(t, []string{testKey2, testKey1}, keys)

	c.Add(testKey3, testValue3)
	keys = c.Keys()
	assert.ElementsMatch(t, []string{testKey3, testKey1}, keys)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_Len(t *testing.T) {
	c, err := NewLRU(2, nil)
	assert.Nil(t, err)
	assert.EqualValues(t, c.Len(), 0)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"
	testKey3 := "test_key_3"
	testValue3 := "test_value_3"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)
	assert.EqualValues(t, c.Len(), 2)

	c.Add(testKey3, testValue3)
	assert.EqualValues(t, c.Len(), 2)
}

func TestLRU_Capacity(t *testing.T) {
	c, err := NewLRU(5, nil)
	assert.Nil(t, err)
	assert.EqualValues(t, c.Len(), 0)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"
	testKey3 := "test_key_3"
	testValue3 := "test_value_3"

	c.Add(testKey1, testValue1)
	assert.EqualValues(t, c.Capacity(), 5)
	c.Add(testKey2, testValue2)
	assert.EqualValues(t, c.Capacity(), 5)

	c.Add(testKey3, testValue3)
	assert.EqualValues(t, c.Capacity(), 5)
}

func TestLRU_Purge(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)
	assert.EqualValues(t, c.Len(), 0)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"
	testKey3 := "test_key_3"
	testValue3 := "test_value_3"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)
	assert.EqualValues(t, c.Len(), 2)

	c.Add(testKey3, testValue3)
	assert.EqualValues(t, c.Len(), 2)

	c.Purge()
	assert.EqualValues(t, c.Len(), 0)
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 3
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_Remove(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)
	assert.EqualValues(t, c.Len(), 0)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)
	assert.EqualValues(t, c.Len(), 2)

	c.Remove(testKey1)
	c.Remove(testKey2)

	assert.EqualValues(t, c.Len(), 0)
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 2
	}, 1*time.Second, 100*time.Millisecond)
}

func TestLRU_RemoveOldest(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)
	assert.EqualValues(t, c.Len(), 0)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)
	assert.EqualValues(t, c.Len(), 2)

	v, ok := c.Get(testKey1)
	assert.True(t, ok)
	assert.EqualValues(t, v, testValue1)

	v, ok = c.Get(testKey2)
	assert.True(t, ok)
	assert.EqualValues(t, v, testValue2)

	c.Remove(testKey1)
	c.Remove(testKey2)

	v, ok = c.Get(testKey1)
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok = c.Get(testKey2)
	assert.False(t, ok)
	assert.Nil(t, v)

	assert.EqualValues(t, c.Len(), 0)
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 2
	}, 1*time.Second, 100*time.Millisecond)

}

func TestLRU_Resize(t *testing.T) {
	evicted := int32(0)
	c, err := NewLRU(2, func(Key, Value) { atomic.AddInt32(&evicted, 1) })
	assert.Nil(t, err)
	assert.EqualValues(t, c.Len(), 0)

	testKey1 := "test_key_1"
	testValue1 := "test_value_1"
	testKey2 := "test_key_2"
	testValue2 := "test_value_2"

	c.Add(testKey1, testValue1)
	c.Add(testKey2, testValue2)
	assert.EqualValues(t, c.Len(), 2)

	v, ok := c.Get(testKey1)
	assert.True(t, ok)
	assert.EqualValues(t, v, testValue1)

	v, ok = c.Get(testKey2)
	assert.True(t, ok)
	assert.EqualValues(t, v, testValue2)

	c.Resize(1)

	v, ok = c.Get(testKey1)
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok = c.Get(testKey2)
	assert.True(t, ok)
	assert.EqualValues(t, v, testValue2)

	assert.EqualValues(t, c.Len(), 1)
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)

	c.Resize(3)

	assert.EqualValues(t, c.Len(), 1)
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&evicted) == 1
	}, 1*time.Second, 100*time.Millisecond)
}
