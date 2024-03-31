package jws

import (
	"encoding/base64"
	"encoding/json"
)

type Header struct {
	alg    Algorithm
	values map[string]any
}

func NewHeader(alg Algorithm) *Header {
	return &Header{
		alg:    alg,
		values: map[string]any{},
	}
}

func (h *Header) SetAlg(alg Algorithm) *Header {
	h.alg = alg
	return h
}

func (h *Header) Set(key string, value any) *Header {
	h.values[key] = value
	return h
}

func (h *Header) Get(key string) any {
	return h.values[key]
}

func (h *Header) Base64Url() (string, error) {
	data, err := json.Marshal(h.values)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
