package platform

import "context"

type CDNItem struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

// CDNManager cdn管理器
type CDNManager interface {
	// GetCDNList 获取cdn列表
	GetCDNList(ctx context.Context) ([]*CDNItem, error)
	// GetCDNDetail 获取cdn详情
	GetCDNDetail(ctx context.Context, domain *CDNItem) (*CDNItem, error)
	// SetCDNHttpsCert 设置cdn https证书
	SetCDNHttpsCert(ctx context.Context, domain *CDNItem, cert, key []byte) error
	// UserDetails 获取用户信息
	UserDetails(ctx context.Context) (*User, error)
}
