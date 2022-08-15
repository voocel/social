package tcp

// PBCodec implements the Codec interface
type PBCodec struct{}

func (c *PBCodec) Encode(v interface{}) ([]byte, error) {
	panic("implement me")
}

func (c *PBCodec) Decode(data []byte, v interface{}) error {
	panic("implement me")
}
