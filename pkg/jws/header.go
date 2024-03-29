package jws

import (
	"encoding/base64"
	"encoding/json"
)

type Header struct {
	alg    Algorithm
	values map[string]interface{}
}

func NewHeader(alg Algorithm) *Header {
	return &Header{
		alg:    alg,
		values: map[string]interface{}{},
	}
}

func (h *Header) Set(key string, value interface{}) *Header {
	h.values[key] = value
	return h
}

func (h *Header) Get(key string) interface{} {
	return h.values[key]
}

func (h *Header) Base64Url() (string, error) {
	data, err := json.Marshal(h.values)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
