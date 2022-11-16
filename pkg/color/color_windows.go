package color

import (
	"fmt"
	"math/rand"
	"strconv"
)

func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

func Yellow(msg string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", msg)
}

func Red(msg string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg)
}

func Redf(msg string, arg interface{}) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m %+v\n", msg, arg)
}

func Blue(msg string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", msg)
}

func Green(msg string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", msg)
}

func Greenf(msg string, arg interface{}) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m %+v\n", msg, arg)
}
