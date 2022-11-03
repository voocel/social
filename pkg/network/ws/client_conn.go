package ws

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
	"social/pkg/network"
	"sync/atomic"
)

type clientConn struct {
	cid    int64
	uid    int64
	client *client
	conn   *websocket.Conn
	state  int32
	sendCh chan []byte
	exitCh chan struct{}
}

func NewWsConn(c *client, conn *websocket.Conn) network.Conn {
	cc := &clientConn{
		cid:    0,
		uid:    0,
		conn:   conn,
		client: c,
		sendCh: make(chan []byte, 1024),
	}
	if cc.client.connectHandler != nil {
		cc.client.connectHandler(cc)
	}
	go cc.readLoop()
	go cc.writeLoop()
	return cc
}

func (cc *clientConn) readLoop() {
	for {
		select {
		case <-cc.exitCh:
			return
		default:
			fmt.Println("read")
			_, rawMsg, err := cc.conn.ReadMessage()
			if err != nil {
				fmt.Println("read err: ", err)
				cc.conn.Close()
				return
			}
			if cc.client.receiveHandler != nil {
				cc.client.receiveHandler(cc, rawMsg, 1)
			}
		}
	}
}

func (cc *clientConn) writeLoop() {
	for {
		select {
		case <-cc.exitCh:
			return
		case msg, ok := <-cc.sendCh:
			if !ok {
				return
			}
			if err := cc.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("client write closed")
				return
			}
		}
	}
}

func (cc *clientConn) checkState() error {
	switch network.ConnState(atomic.LoadInt32(&cc.state)) {
	case network.ConnHanged:
		return network.ErrConnectionHanged
	case network.ConnClosed:
		return network.ErrConnectionClosed
	}
	return nil
}

func (cc *clientConn) Cid() int64 {
	return cc.cid
}

func (cc *clientConn) Uid() int64 {
	return cc.uid
}

func (cc *clientConn) Bind(uid int64) {
	cc.uid = uid
}

func (cc *clientConn) Send(msg []byte, msgType ...int) error {
	if err := cc.checkState(); err != nil {
		return err
	}
	return cc.conn.WriteMessage(websocket.TextMessage, msg)
}

func (cc *clientConn) SendMessage() {

}

func (cc *clientConn) AsyncSend(msg []byte, msgType ...int) error {
	if err := cc.checkState(); err != nil {
		return err
	}
	cc.sendCh <- msg
	return nil
}

func (cc *clientConn) State() network.ConnState {
	return network.ConnState(atomic.LoadInt32(&cc.state))
}

func (cc *clientConn) Close() error {
	if err := cc.checkState(); err != nil {
		return err
	}
	atomic.StoreInt32(&cc.state, int32(network.ConnHanged))
	close(cc.sendCh)
	return cc.conn.Close()
}

func (cc *clientConn) LocalIP() string {
	return ExtractIP(cc.LocalAddr())
}

func (cc *clientConn) LocalAddr() net.Addr {
	if err := cc.checkState(); err != nil {
		return nil
	}

	return cc.conn.LocalAddr()
}

func (cc *clientConn) RemoteIP() string {
	return ExtractIP(cc.RemoteAddr())
}

func (cc *clientConn) RemoteAddr() net.Addr {
	if err := cc.checkState(); err != nil {
		return nil
	}

	return cc.conn.RemoteAddr()
}

func (cc *clientConn) close() {
	atomic.StoreInt32(&cc.state, int32(network.ConnClosed))

	if cc.client.disconnectHandler != nil {
		cc.client.disconnectHandler(cc, nil)
	}
}

func ExtractIP(addr net.Addr) (host string) {
	host, _, _ = net.SplitHostPort(addr.String())
	return
}
