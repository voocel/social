package packet

import (
	"bytes"
	"encoding/binary"
)

const (
	seqLenBytes   = 2
	routeLenBytes = 4
)

// Pack 打包消息
func Pack(message *Message) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(len(message.Buffer) + seqLenBytes + routeLenBytes)

	if err := binary.Write(&buf, binary.LittleEndian, int16(message.Seq)); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.LittleEndian, message.Route); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.LittleEndian, message.Buffer); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unpack 解包消息
func Unpack(data []byte) (*Message, error) {
	reader := bytes.NewReader(data)
	packet := &Message{Buffer: make([]byte, reader.Len()-seqLenBytes-routeLenBytes)}

	var seq int16
	if err := binary.Read(reader, binary.LittleEndian, &seq); err != nil {
		return nil, err
	}
	packet.Seq = int32(seq)

	if err := binary.Read(reader, binary.LittleEndian, &packet.Route); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &packet.Buffer); err != nil {
		return nil, err
	}

	return packet, nil
}
