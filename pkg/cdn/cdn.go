package cdn

import "context"

// CDN cdn接口
type CDN interface {
	// GetDomainList 获取域名列表
	GetDomainList(ctx context.Context) ([]*Domain, error)
	// SetHttpsCert 设置https证书
	SetHttpsCert(ctx context.Context, domain string, cert, key []byte) error
}
