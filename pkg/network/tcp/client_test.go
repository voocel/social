package tcp

import (
	"fmt"
	"social/pkg/message"
	"social/pkg/network"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient("127.0.0.1:8899")

	c.OnConnect(func(conn network.Conn) {
		fmt.Println("connect")
	})

	c.OnReceive(func(conn network.Conn, msg *message.Message, msgType int) {
		fmt.Println("msg: ", msg)
	})

	c.OnDisconnect(func(conn network.Conn, err error) {
		fmt.Println("disconnect: ", err)
	})

	conn, err := c.Dial()
	if err != nil {
		panic(err)
	}

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
}
