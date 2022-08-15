package im

import (
	"os"
	"os/signal"
	"social/pkg/log"
	"social/pkg/websocket"
	"syscall"
)

func Run()  {
	ws := websocket.NewWsServer()
	go ws.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		s := <-ch
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Sync()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}