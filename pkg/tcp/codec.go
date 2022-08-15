package tcp

type Codec interface {
	// Encode encodes data into []byte
	// Returns error when error occurred
	Encode(v interface{}) ([]byte, error)

	// Decode decodes data into v
	// Returns error when error occurred
	Decode(data []byte, v interface{}) error
}
