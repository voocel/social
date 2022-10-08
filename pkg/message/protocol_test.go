package message

import (
	"bytes"
	"testing"
)

func TestProtocol(t *testing.T) {
	msg := NewMessage(Heartbeat, []byte("hello123"))
	p := NewDefaultProtocol()
	b, err := p.Pack(msg)
	t.Log(b, err)

	s := []byte("aa")
	m, err := p.Unpack(bytes.NewReader(s))
	t.Log(m, err)
}
