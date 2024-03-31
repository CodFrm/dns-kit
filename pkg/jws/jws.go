package jws

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// https://datatracker.ietf.org/doc/html/rfc7515

func Encode(header *Header, payload any, opts ...EncodeOption) (string, error) {
	options := newEncodeOptions(opts...)
	if err := header.alg.PreCompute(header); err != nil {
		return "", err
	}
	headerStr, err := header.Base64Url()
	if err != nil {
		return "", err
	}
	var payloadStr string
	switch p := payload.(type) {
	case string:
		payloadStr = base64.RawURLEncoding.EncodeToString([]byte(p))
	case []byte:
		payloadStr = base64.RawURLEncoding.EncodeToString(p)
	case nil:
		payloadStr = ""
	default:
		payloadData, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}
		payloadStr = base64.RawURLEncoding.EncodeToString(payloadData)
	}
	signatureData, err := header.alg.Compute([]byte(headerStr + "." + payloadStr))
	if err != nil {
		return "", err
	}
	signatureStr := base64.RawURLEncoding.EncodeToString(signatureData)
	return options.serialization(headerStr, payloadStr, signatureStr), nil
}

var ErrInvalidSignature = errors.New("invalid signature")

func Decode(data string, header *Header, payload any, opts ...DecodeOption) error {
	options := newDecodeOptions(opts...)
	headerStr, payloadStr, signatureStr, err := options.unmarshaler(data)
	if err != nil {
		return err
	}
	headerData, err := base64.RawURLEncoding.DecodeString(headerStr)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(headerData, &header.values); err != nil {
		return err
	}
	if err := header.alg.PreVerify(header); err != nil {
		return err
	}
	// 验证签名
	signature, err := base64.RawURLEncoding.DecodeString(signatureStr)
	if err != nil {
		return err
	}
	err = header.alg.Verify([]byte(headerStr+"."+payloadStr), signature)
	if err != nil {
		return err
	}
	payloadData, err := base64.RawURLEncoding.DecodeString(payloadStr)
	if err != nil {
		return err
	}
	switch p := payload.(type) {
	case *string:
		*p = string(payloadData)
	case *[]byte:
		*p = payloadData
	default:
		if err := json.Unmarshal(payloadData, &payload); err != nil {
			return err
		}
	}
	return nil
}
