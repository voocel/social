package util

import (
	"fmt"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var letterRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return b
}

func RandomString(n int) string {
	return string(RandomBytes(n))
}

func GenRandomBytes() []byte {
	l := rand.Uint64()%10 + 1
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return nil
	}
	return b
}

func GenRandomStr() string {
	return fmt.Sprintf("%X", GenRandomBytes())
}
