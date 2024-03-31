package acme

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"math/big"

	"github.com/codfrm/dns-kit/pkg/jws"
)

type es256 struct {
	jws.Algorithm
	kid        string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	jwk        map[string]any
}

func newEs256(kid string, privateKey *ecdsa.PrivateKey) *es256 {
	return &es256{
		Algorithm:  jws.ES256(privateKey),
		kid:        kid,
		privateKey: privateKey,
	}
}

func (h *es256) PreCompute(header *jws.Header) error {
	header.Set("alg", "ES256").
		Set("kid", h.kid)
	return nil
}

func (h *es256) PreVerify(header *jws.Header) error {
	h.publicKey = &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     nil,
		Y:     nil,
	}
	// 解析x y
	x, err := base64.RawURLEncoding.DecodeString(h.jwk["x"].(string))
	if err != nil {
		return err
	}
	y, err := base64.RawURLEncoding.DecodeString(h.jwk["y"].(string))
	if err != nil {
		return err
	}
	h.publicKey.X = new(big.Int).SetBytes(x)
	h.publicKey.Y = new(big.Int).SetBytes(y)
	return nil
}

func (h *es256) Verify(data []byte, signature []byte) error {
	// 对消息进行SHA-256哈希处理
	hash := sha256.New()
	hash.Write(data)

	// 解析签名
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])

	// 验证签名
	if !ecdsa.Verify(h.publicKey, hash.Sum(nil), r, s) {
		return jws.ErrInvalidSignature
	}
	return nil
}
