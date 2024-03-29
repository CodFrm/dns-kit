package acme

import (
	"crypto/ecdsa"

	"github.com/codfrm/dns-kit/pkg/jws"
)

type es256 struct {
	jws.Algorithm
	kid        string
	privateKey *ecdsa.PrivateKey
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
