package ws

import (
	"fmt"
	"social/pkg/message"
	"social/pkg/network"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient("ws://127.0.0.1:8688/ws")
	c.OnConnect(func(conn network.Conn) {
		fmt.Println("client connect successfully")
	})
	c.OnReceive(func(conn network.Conn, msg *message.Message, msgType int) {
		fmt.Printf("收到服务端消息: %+v\n", msg)
	})
	c.OnDisconnect(func(conn network.Conn, err error) {
		fmt.Println("关闭连接: ", conn.Cid())
	})
	conn, err := c.Dial()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		msg := message.NewMessage(message.Heartbeat, []byte(fmt.Sprintf("test-%d", i)))
		b, err := c.Codec().Pack(msg)
		if err != nil {
			fmt.Println("pack err:", err)
		}
		if err = conn.Send(b, 1); err != nil {
			fmt.Println("send err:", err)
		}
	}

	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	//for {
	//	s := <-ch
	//	switch s {
	//	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
	//		log.Sync()
	//		return
	//	case syscall.SIGHUP:
	//	default:
	//		return
	//	}
	//}
}
