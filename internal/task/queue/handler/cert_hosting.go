package handler

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	logger2 "github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
	"github.com/codfrm/dns-kit/internal/repository/cdn_repo"
	"github.com/codfrm/dns-kit/internal/repository/cert_hosting_repo"
	"github.com/codfrm/dns-kit/internal/repository/cert_repo"
	"github.com/codfrm/dns-kit/internal/task/queue"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
	"github.com/codfrm/dns-kit/pkg/platform"
	"go.uber.org/zap"
)

type CertHostingHandler struct {
}

func (h *CertHostingHandler) Register(ctx context.Context) error {
	if err := queue.SubscribeCertCreateAfter(ctx, h.CreateCertAfter, broker.Group("cert_hosting")); err != nil {
		return err
	}
	return nil
}

// CreateCertAfter 证书创建完成后, 查询有没有关联的托管，有的部署到cdn
func (h *CertHostingHandler) CreateCertAfter(ctx context.Context, msg *message.CreateCertAfterMessage) error {
	list, err := cert_hosting_repo.CertHosting().FindByCert(ctx, msg.ID)
	if err != nil {
		return err
	}
	cert, err := cert_repo.Cert().Find(ctx, msg.ID)
	if err != nil {
		return err
	}
	// 一个证书对应多个托管, 为以后的泛域名做准备?
	for _, v := range list {
		logger := logger2.Ctx(ctx).With(zap.Int64("cert_id", msg.ID), zap.Int64("hosting_id", v.ID))
		if !msg.Success {
			// 更新为部署失败
			if err := cert_hosting_repo.CertHosting().UpdateStatus(ctx, v.ID, cert_hosting_entity.CertHostingStatusDeployFail); err != nil {
				logger.Error("update hosting status failed", zap.Error(err))
			}
			return nil
		}
		// 部署到cdn
		err := func() error {
			cdn, err := cdn_repo.Cdn().Find(ctx, v.CdnID)
			if err != nil {
				logger.Error("find cdn error", zap.Error(err))
				return err
			}
			if err := cdn.Check(ctx); err != nil {
				logger.Error("cdn check error", zap.Error(err))
				return err
			}
			manager, err := cdn.CDNManger(ctx)
			if err != nil {
				logger.Error("cdn manager error", zap.Error(err))
				return err
			}
			if err := manager.SetCDNHttpsCert(ctx, &platform.CDNItem{
				ID:     cdn.CdnID,
				Domain: cdn.Domain,
			}, cert.Certificate, cert.PrivateKey); err != nil {
				logger.Error("set cdn https cert error", zap.Error(err))
				return err
			}
			return nil
		}()
		if err != nil {
			// 修改状态
			if err := cert_hosting_repo.CertHosting().UpdateStatus(ctx, v.ID, cert_hosting_entity.CertHostingStatusDeployFail); err != nil {
				logger.Error("update hosting status failed", zap.Error(err))
			}
		} else {
			// 更新状态并更新
			v.CertID = msg.ID
			v.Status = cert_hosting_entity.CertHostingStatusActive
			if err := cert_hosting_repo.CertHosting().Update(ctx, v); err != nil {
				logger.Error("update hosting status failed", zap.Error(err))
			} else {
				logger.Info("deploy success")
			}
		}
	}
	return nil
}
