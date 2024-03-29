package jws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	data, err := Encode(NewHeader(HS256([]byte("123456"))),
		`{"iss":"joe","exp":1300819380,"http://example.com/is_root":true}`)
	assert.NoError(t, err)
	payload := map[string]interface{}{}
	header := NewHeader(HS256([]byte("123456")))
	err = Decode(data, header, &payload)
	assert.NoError(t, err)
	assert.Equal(t, "joe", payload["iss"])
	assert.Equal(t, float64(1300819380), payload["exp"])
	assert.Equal(t, true, payload["http://example.com/is_root"])
	assert.Equal(t, "HS256", header.Get("alg"))
}
