package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"golang.org/x/net/publicsuffix"
)

func GetTLDMap(domains []string) (map[string]struct{}, error) {
	// 获取顶级域名
	domainMap := make(map[string]struct{})
	for _, v := range domains {
		tld, err := publicsuffix.EffectiveTLDPlusOne(v)
		if err != nil {
			return nil, err
		}
		domainMap[tld] = struct{}{}
	}
	return domainMap, nil
}

func GetTLD(domain string) (string, error) {
	return publicsuffix.EffectiveTLDPlusOne(domain)
}

func DecodeCertPEM(certPem string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPem))
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
