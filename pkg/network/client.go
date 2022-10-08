package network

import "social/pkg/message"

type Client interface {
	Dial() (Conn, error)

	Codec() (p message.Protocol)

	OnConnect(handler ConnectHandler)

	OnReceive(handler ReceiveHandler)

	OnDisconnect(handler DisconnectHandler)
}
