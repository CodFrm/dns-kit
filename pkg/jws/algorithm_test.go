package jws

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestES256(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal("generate key faild")
	}
	data, err := Encode(NewHeader(ES256(privateKey)),
		`{"iss":"joe","exp":1300819380,"http://example.com/is_root":true}`)
	assert.NoError(t, err)
	payload := map[string]interface{}{}
	header := NewHeader(ES256(nil))
	err = Decode(data, header, &payload)
	assert.NoError(t, err)
	assert.Equal(t, "joe", payload["iss"])
	assert.Equal(t, float64(1300819380), payload["exp"])
	assert.Equal(t, true, payload["http://example.com/is_root"])
	assert.Equal(t, "ES256", header.Get("alg"))
}
