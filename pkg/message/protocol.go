package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Won't compile if Protocol can't be realized by a DefaultProtocol
var _ Protocol = &DefaultProtocol{}

type Protocol interface {
	// Pack packs Message into the packet to be written
	Pack(msg *Message) ([]byte, error)

	// Unpack unpacks the message packet from reader
	Unpack(reader io.Reader) (*Message, error)
}

// DefaultProtocol is the default packet
// ╔═══════════╤════════╤═════════╗
// ║ FIELD     │ TYPE   │  SIZE   ║
// ╠═══════════╪════════╪═════════╣
// ║ Cmd       │ uint16 │ 2       ║
// ║ Size      │ uint32 │ 4       ║
// ║ Checksum  │ uint32 │ 4       ║
// ║ Data      │ []byte │ dynamic ║
// ╚═══════════╧════════╧═════════╝
type DefaultProtocol struct{}

// NewDefaultProtocol create a *DefaultPacker with initial field value.
func NewDefaultProtocol() *DefaultProtocol {
	return &DefaultProtocol{}
}
// 
// Pack encodes the message into bytes data
func (p *DefaultProtocol) Pack(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, msg.size)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, msg.cmd)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, msg.data)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, msg.checksum)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unpack decodes the io.Reader into Message
func (p *DefaultProtocol) Unpack(reader io.Reader) (msg *Message, err error) {
	msg = &Message{}
	err = binary.Read(reader, binary.LittleEndian, &msg.size)
	if err != nil {
		return
	}

	err = binary.Read(reader, binary.LittleEndian, &msg.cmd)
	if err != nil {
		return
	}

	msg.data = make([]byte, msg.size)
	err = binary.Read(reader, binary.LittleEndian, &msg.data)
	if err != nil {
		return
	}

	err = binary.Read(reader, binary.LittleEndian, &msg.checksum)
	if err != nil {
		return
	}

	if !msg.Checksum() {
		return nil, errors.New("checksum error")
	}

	return
}
