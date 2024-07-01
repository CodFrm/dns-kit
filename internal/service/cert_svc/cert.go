package cert_svc

import (
	"context"
	"sync"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/dns-kit/internal/model/entity/acme_entity"
	"github.com/codfrm/dns-kit/internal/repository/acme_repo"
	"github.com/codfrm/dns-kit/pkg/acme"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/pkg/utils"
	"github.com/codfrm/dns-kit/internal/repository/cert_repo"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"
	"github.com/codfrm/dns-kit/internal/task/queue"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
	"go.uber.org/zap"

	api "github.com/codfrm/dns-kit/internal/api/cert"
)

type CertSvc interface {
	// List 获取证书列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Create 创建证书
	Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error)
	// NewACME 创建ACME
	NewACME(ctx context.Context, email string) (*acme.Acme, error)
	// Download 下载证书
	Download(ctx context.Context, req *api.DownloadRequest) (*api.DownloadResponse, error)
	// Delete 删除证书
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
	// CheckDomains 检查域名
	CheckDomains(ctx context.Context, domains []string) error
}

type certSvc struct {
	sync.Mutex
}

var defaultCert = &certSvc{}

func Cert() CertSvc {
	return defaultCert
}

// List 获取证书列表
func (c *certSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, total, err := cert_repo.Cert().FindPage(ctx, req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Item]{
			List:  make([]*api.Item, 0),
			Total: total,
		},
	}
	for _, v := range list {
		resp.PageResponse.List = append(resp.PageResponse.List, &api.Item{
			ID:         v.ID,
			Email:      v.Email,
			Domains:    v.Domains,
			Status:     v.Status,
			Createtime: v.Createtime,
			Expiretime: v.Expiretime,
		})
	}
	return resp, nil
}

func (c *certSvc) CheckDomains(ctx context.Context, domains []string) error {
	// 获取顶级域名
	domainMap, err := utils.GetTLDMap(domains)
	if err != nil {
		return i18n.NewError(ctx, code.InvalidDomain)
	}
	// 搜索域名是否在纳管中
	for domain := range domainMap {
		entity, err := domain_repo.Domain().FindByDomain(ctx, domain)
		if err != nil {
			return err
		}
		if err := entity.Check(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Create 创建证书
func (c *certSvc) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	if err := c.CheckDomains(ctx, req.Domains); err != nil {
		return nil, err
	}
	cert := &cert_entity.Cert{
		Email:      req.Email,
		Domains:    req.Domains,
		Status:     cert_entity.CertStatusApply,
		Createtime: time.Now().Unix(),
	}
	if err := cert_repo.Cert().Create(ctx, cert); err != nil {
		return nil, err
	}
	if !req.GetIgnoreMsg() {
		if err := queue.PublishCertCreate(ctx, &message.CreateCertMessage{ID: cert.ID}); err != nil {
			return nil, err
		}
	}
	_ = audit.Ctx(ctx).Record("create", zap.Int64("id", cert.ID),
		zap.Strings("domains", req.Domains), zap.String("email", cert.Email))
	return &api.CreateResponse{ID: cert.ID}, nil
}

func (c *certSvc) NewACME(ctx context.Context, email string) (*acme.Acme, error) {
	c.Lock()
	defer c.Unlock()
	acmeEntity, err := acme_repo.Acme().FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if acmeEntity != nil {
		// 存在利用之前的
		return acmeEntity.NewACME()
	}
	privateKey, err := acme.GenerateKey()
	if err != nil {
		return nil, err
	}
	acmeEntity = &acme_entity.Acme{
		Email:      email,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := acmeEntity.SavePrivateKey(privateKey); err != nil {
		return nil, err
	}
	acmeInstance, err := acmeEntity.NewACME()
	if err != nil {
		return nil, err
	}
	kid, err := acmeInstance.NewAccount(ctx)
	if err != nil {
		return nil, err
	}
	acmeEntity.Kid = kid
	if err := acme_repo.Acme().Create(ctx, acmeEntity); err != nil {
		return nil, err
	}
	return acmeInstance, nil
}

// Download 下载证书
func (c *certSvc) Download(ctx context.Context, req *api.DownloadRequest) (*api.DownloadResponse, error) {
	cert, err := cert_repo.Cert().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := cert.Check(ctx); err != nil {
		return nil, err
	}
	if cert.Status != cert_entity.CertStatusActive {
		return nil, i18n.NewError(ctx, code.CertNotActive)
	}
	_ = audit.Ctx(ctx).Record("download", zap.Int64("id", req.ID))
	return &api.DownloadResponse{Cert: cert.Certificate, CSR: cert.CertificateRequest, Key: cert.PrivateKey}, nil
}

// Delete 删除证书
func (c *certSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	cert, err := cert_repo.Cert().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := cert.Check(ctx); err != nil {
		return nil, err
	}
	if cert.Status == cert_entity.CertStatusApply {
		return nil, i18n.NewError(ctx, code.CertStatusApply)
	}
	if err := cert_repo.Cert().Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("delete", zap.Int64("id", req.ID))
	return nil, nil
}
