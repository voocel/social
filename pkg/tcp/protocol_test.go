package tcp

import (
	"bytes"
	"testing"
)

func TestDefaultProtocol(t *testing.T) {
	m := NewMessage(Heartbeat, []byte("ping"))
	p := NewDefaultProtocol()

	b, err := p.Pack(m)
	if err != nil {
		t.Fatal(err)
	}

	res, err := p.Unpack(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	t.Log(res)
}
