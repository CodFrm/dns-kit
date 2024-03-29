package jws

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

type Algorithm interface {
	PreCompute(header *Header) error
	Compute(data []byte) ([]byte, error)
	PreVerify(header *Header) error
	Verify(data, signature []byte) error
}

type hs256 struct {
	secret []byte
}

func HS256(secret []byte) Algorithm {
	return &hs256{secret}
}

func (h *hs256) PreCompute(header *Header) error {
	header.Set("alg", "HS256")
	return nil
}

func (h *hs256) Compute(data []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, h.secret)
	hash.Write(data)
	return hash.Sum(nil), nil
}

func (h *hs256) PreVerify(header *Header) error {
	return nil
}

func (h *hs256) Verify(data, signature []byte) error {
	hash := hmac.New(sha256.New, h.secret)
	hash.Write(data)
	if !hmac.Equal(hash.Sum(nil), signature) {
		return ErrInvalidSignature
	}
	return nil
}

type es256 struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func ES256(privateKey *ecdsa.PrivateKey) Algorithm {
	ret := &es256{
		privateKey: privateKey,
	}
	return ret
}

func ES256Jwk(publicKey ecdsa.PublicKey) string {
	return fmt.Sprintf(`{"crv":"%s","kty":"EC","x":"%s","y":"%s"}`,
		"P-256",
		base64.RawURLEncoding.EncodeToString(publicKey.X.Bytes()),
		base64.RawURLEncoding.EncodeToString(publicKey.Y.Bytes()),
	)
}

func (e *es256) PreCompute(header *Header) error {
	header.Set("alg", "ES256").
		Set("jwk", json.RawMessage(ES256Jwk(e.privateKey.PublicKey)))
	return nil
}

// Compute 加密 参考: https://datatracker.ietf.org/doc/html/rfc7515#appendix-A.3
func (e *es256) Compute(data []byte) ([]byte, error) {
	// 对消息进行SHA-256哈希处理
	hash := sha256.New()
	hash.Write(data)

	// 使用私钥对哈希进行签名
	r, s, err := ecdsa.Sign(rand.Reader, e.privateKey, hash.Sum(nil))
	if err != nil {
		return nil, err
	}

	rb, sb := r.Bytes(), s.Bytes()
	size := e.privateKey.PublicKey.Params().BitSize / 8
	if size%8 > 0 {
		size++
	}
	sig := make([]byte, size*2)
	copy(sig[size-len(rb):], rb)
	copy(sig[size*2-len(sb):], sb)

	return sig, nil
}

func (e *es256) PreVerify(header *Header) error {
	e.publicKey = &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     nil,
		Y:     nil,
	}
	// 解析x y
	jwk, ok := header.Get("jwk").(map[string]interface{})
	if !ok {
		return errors.New("jwk not found")
	}
	x, err := base64.RawURLEncoding.DecodeString(jwk["x"].(string))
	if err != nil {
		return err
	}
	y, err := base64.RawURLEncoding.DecodeString(jwk["y"].(string))
	if err != nil {
		return err
	}
	e.publicKey.X = new(big.Int).SetBytes(x)
	e.publicKey.Y = new(big.Int).SetBytes(y)
	return nil
}

func (e *es256) Verify(data []byte, signature []byte) error {
	// 对消息进行SHA-256哈希处理
	hash := sha256.New()
	hash.Write(data)

	// 解析签名
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])

	// 验证签名
	if !ecdsa.Verify(e.publicKey, hash.Sum(nil), r, s) {
		return ErrInvalidSignature
	}
	return nil
}
