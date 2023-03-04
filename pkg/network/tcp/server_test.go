package tcp

import (
	"fmt"
	"os"
	"os/signal"
	"social/pkg/log"
	"social/pkg/network"
	"syscall"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(":8899")

	s.OnStart(func() {
		fmt.Println("start")
	})

	s.OnConnect(func(conn network.Conn) {
		fmt.Println("connect")
	})

	s.OnReceive(func(conn network.Conn, msg []byte) {
		fmt.Println("msg: ", string(msg))
	})

	s.OnDisconnect(func(conn network.Conn, err error) {
		fmt.Println("disconnect: ", err)
	})

	s.OnStop(func() {
		fmt.Println("stop")
	})

	if err := s.Start(); err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Sync()
			s.Stop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
