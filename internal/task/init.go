package task

import (
	"context"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/dns-kit/internal/task/queue"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
	"go.uber.org/zap"

	"github.com/codfrm/cago/pkg/gogo"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/codfrm/dns-kit/internal/repository/cert_repo"
)

// InitTask 初始化的任务
func InitTask() error {
	if err := gogo.Go(RetryApplyCert); err != nil {
		return err
	}
	return nil
}

// RetryApplyCert 重试申请证书
func RetryApplyCert(ctx context.Context) error {
	// 将所有状态为申请中的证书重新运行一次
	list, err := cert_repo.Cert().FindByStatus(ctx, cert_entity.CertStatusApply)
	if err != nil {
		return err
	}
	for _, v := range list {
		if err := queue.PublishCertCreate(ctx, &message.CreateCertMessage{ID: v.ID}); err != nil {
			logger.Ctx(ctx).Error("publish cert create failed", zap.Error(err))
		}
	}
	return nil
}
