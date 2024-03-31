package jws

import (
	"encoding/json"
	"errors"
	"strings"
)

type EncodeOption func(o *EncodeOptions)

type EncodeOptions struct {
	serialization Serialization
}

type DecodeOption func(o *DecodeOptions)

type DecodeOptions struct {
	unmarshaler Unmarshaler
}

func newEncodeOptions(opts ...EncodeOption) EncodeOptions {
	opt := EncodeOptions{
		serialization: CompactSerialization,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func newDecodeOptions(opts ...DecodeOption) DecodeOptions {
	opt := DecodeOptions{
		unmarshaler: CompactUnmarshaler,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithSerialization(s Serialization) EncodeOption {
	return func(o *EncodeOptions) {
		o.serialization = s
	}
}

func WithUnmarshaler(u Unmarshaler) DecodeOption {
	return func(o *DecodeOptions) {
		o.unmarshaler = u
	}
}

type Serialization func(header, payload, signature string) string

type Unmarshaler func(data string) (string, string, string, error)

func CompactSerialization(header, payload, signature string) string {
	return header + "." + payload + "." + signature
}

var ErrInvalidToken = errors.New("invalid token")

func CompactUnmarshaler(data string) (string, string, string, error) {
	parts := strings.Split(data, ".")
	if len(parts) != 3 {
		return "", "", "", ErrInvalidToken
	}
	return parts[0], parts[1], parts[2], nil
}

func JSONSerialization(header, payload, signature string) string {
	m := map[string]string{
		"protected": header,
		"payload":   payload,
		"signature": signature,
	}
	data, _ := json.Marshal(m)
	return string(data)
}

func JSONUnmarshaler(data string) (string, string, string, error) {
	var m map[string]string
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return "", "", "", err
	}
	return m["protected"], m["payload"], m["signature"], nil
}
