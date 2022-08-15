package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
)

type Message struct {
	cmd      CMD
	size     uint32
	data     []byte
	checksum uint32
}

type CMD uint16

const (
	Heartbeat CMD = iota
	Ack
	Single
	All
)

func NewMessage(cmd CMD, data []byte) *Message {
	msg := &Message{
		cmd:  cmd,
		size: uint32(len(data)),
		data: data,
	}
	msg.checksum = msg.calc()
	return msg
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) GetCmd() CMD {
	return m.cmd
}

func (m *Message) GetSize() uint32 {
	return m.size
}

func (m *Message) Checksum() bool {
	return m.checksum == m.calc()
}

func (m *Message) calc() (checksum uint32) {
	if m == nil {
		return
	}

	data := new(bytes.Buffer)
	err := binary.Write(data, binary.LittleEndian, m.cmd)
	if err != nil {
		return
	}
	err = binary.Write(data, binary.LittleEndian, m.data)
	if err != nil {
		return
	}

	checksum = adler32.Checksum(data.Bytes())
	return
}

func (m *Message) String() string {
	return fmt.Sprintf("{cmd:%d, size:%d, data:%v, checksum:%d}", m.GetCmd(), m.GetSize(), string(m.GetData()), m.checksum)
}
