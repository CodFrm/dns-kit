package crontab

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
	"github.com/codfrm/dns-kit/internal/repository/cert_hosting_repo"
	"github.com/codfrm/dns-kit/internal/repository/cert_repo"
	"github.com/codfrm/dns-kit/internal/service/cert_svc"
	"github.com/codfrm/dns-kit/internal/task/queue"
	"github.com/codfrm/dns-kit/internal/task/queue/handler"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
	"go.uber.org/zap"
)

// CheckCertHosting 检查证书托管状态
func CheckCertHosting(ctx context.Context) error {
	hosting, _, err := cert_hosting_repo.CertHosting().FindPage(ctx, httputils.PageRequest{Size: 100})
	if err != nil {
		return err
	}
	for _, v := range hosting {
		logger := logger.Ctx(ctx).With(zap.Int64("cert_id", v.CertID), zap.Int64("hosting_id", v.ID))
		switch v.Status {
		case cert_hosting_entity.CertHostingStatusActive:
			// 检查证书是否过期或者快要过期
			cert, err := cert_repo.Cert().Find(ctx, v.CertID)
			if err != nil {
				logger.Error("find cert error", zap.Error(err))
				continue
			}
			// 如果还有30天到期
			expireTime := time.Unix(cert.Expiretime, 0)
			if time.Now().AddDate(0, 0, 30).After(expireTime) {
				// 更新状态，重新部署
				if err := cert_hosting_repo.CertHosting().UpdateStatus(ctx, v.ID, cert_hosting_entity.CertHostingStatusDeploy); err != nil {
					logger.Error("update hosting status failed", zap.Error(err))
				} else {
					if err := cert_svc.Hosting().ReDeploy(context.Background(), v.ID); err != nil {
						logger.Error("redeploy failed", zap.Error(err))
						continue
					}
				}
			}
		case cert_hosting_entity.CertHostingStatusDeployFail:
			// 重试部署
			cert, err := cert_repo.Cert().Find(ctx, v.CertID)
			if err != nil {
				logger.Error("find cert error", zap.Error(err))
				continue
			}
			if cert == nil {
				// 不存在，更新状态，重新申请部署
				if err := cert_hosting_repo.CertHosting().UpdateStatus(ctx, v.ID, cert_hosting_entity.CertHostingStatusDeploy); err != nil {
					logger.Error("update hosting status failed", zap.Error(err))
				} else {
					if err := cert_svc.Hosting().ReDeploy(context.Background(), v.ID); err != nil {
						logger.Error("redeploy failed", zap.Error(err))
						continue
					}
				}
				continue
			}
			switch cert.Status {
			case cert_entity.CertStatusActive:
				// 如果证书申请成功，重新部署
				if err := (&handler.CertHostingHandler{}).CreateCertAfter(ctx, &message.CreateCertAfterMessage{
					ID:      cert.ID,
					Success: true,
				}); err != nil {
					logger.Error("redeploy failed", zap.Error(err))
				}
			case cert_entity.CertStatusApplyFail, cert_entity.CertStatusExpire:
				// 如果证书申请失败，重新申请证书并部署
				if err := cert_repo.Cert().UpdateStatus(ctx, cert.ID, cert_entity.CertStatusApply); err != nil {
					logger.Error("update cert status failed", zap.Error(err))
				} else {
					if err := queue.PublishCertCreate(ctx, &message.CreateCertMessage{ID: cert.ID}); err != nil {
						logger.Error("redeploy failed", zap.Error(err))
					}
				}
			}
		}
	}
	return nil
}
