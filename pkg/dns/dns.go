package dns

import "context"

// DomainManager 域名管理器
type DomainManager interface {
	// GetDomainList 获取域名列表
	GetDomainList(ctx context.Context) ([]*Domain, error)
	// BuildDNSManager 构建dns域名管理器
	BuildDNSManager(ctx context.Context, domain *Domain) (Manager, error)
	// UserDetails 获取用户信息
	UserDetails(ctx context.Context) (*User, error)
}

// Manager dns管理器
type Manager interface {
	// GetRecordList 获取记录列表 获取所有记录
	GetRecordList(ctx context.Context) ([]*Record, error)
	// AddRecord 添加记录
	AddRecord(ctx context.Context, record *Record) error
	// UpdateRecord 更新记录
	UpdateRecord(ctx context.Context, recordId string, record *Record) error
	// DelRecord 删除记录
	DelRecord(ctx context.Context, recordId string) error
	// ExtraFields 额外字段定义描述
	ExtraFields() []*Extra
}
