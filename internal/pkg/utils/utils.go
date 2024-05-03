package utils

import "golang.org/x/net/publicsuffix"

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
