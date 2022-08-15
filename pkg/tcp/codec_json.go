package tcp

import "encoding/json"

var _ Codec = &JSONCodec{}

// JSONCodec implements the Codec interface
type JSONCodec struct{}

// Encode implements the Codec Encode method
func (c *JSONCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode implements the Codec Decode method
func (c *JSONCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
