package im

import (
	"os"
	"os/signal"
	"social/internal/node"
	"social/pkg/log"
	"syscall"
)

func Run() {
	n := node.NewNode()
	core := NewCore(n.GetProxy())
	core.Init()
	n.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			n.Stop()
			log.Sync()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
