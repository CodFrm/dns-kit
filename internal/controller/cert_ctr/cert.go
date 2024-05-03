package cert_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/cert"
	"github.com/codfrm/dns-kit/internal/service/cert_svc"
)

type Cert struct {
}

func NewCert() *Cert {
	return &Cert{}
}

// List 获取证书列表
func (c *Cert) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return cert_svc.Cert().List(ctx, req)
}

// Create 创建证书
func (c *Cert) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	return cert_svc.Cert().Create(ctx, req)
}
