package acme

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"golang.org/x/crypto/acme/autocert"
)

func CreateCertificateRequest(auth []Identifier) ([]byte, *ecdsa.PrivateKey, error) {
	csr := x509.CertificateRequest{}
	for _, v := range auth {
		switch v.Type {
		case "dns":
			csr.DNSNames = append(csr.DNSNames, v.Value)
		}
	}
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	b, err := x509.CreateCertificateRequest(rand.Reader, &csr, k)
	if err != nil {
		return nil, nil, err
	}
	autocert.NewListener()
	return b, k, nil
}

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
