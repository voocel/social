package tcp

import (
	"google.golang.org/protobuf/proto"
)

// PBCodec implements the Codec interface
type PBCodec struct{}

func (c *PBCodec) Encode(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (c *PBCodec) Decode(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
