package tcp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONCodec_Encode(t *testing.T) {
	c := &JSONCodec{}
	v := struct {
		Id int `json:"id"`
	}{Id: 1}
	b, err := c.Encode(v)
	assert.NoError(t, err)
	assert.JSONEq(t, string(b), `{"id": 1}`)
}

func TestJSONCodec_Decode(t *testing.T) {
	c := &JSONCodec{}
	data := []byte(`{"id": 1}`)
	var v struct {
		Id int `json:"id"`
	}
	err := c.Decode(data, &v)
	assert.NoError(t, err)
	assert.EqualValues(t, v.Id, 1)
}
