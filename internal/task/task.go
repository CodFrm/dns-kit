package task

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/cron"
	"github.com/codfrm/dns-kit/internal/task/crontab"
	"github.com/codfrm/dns-kit/internal/task/queue/handler"
)

type Handler interface {
	Register(ctx context.Context) error
}

func Task(ctx context.Context, config *configs.Config) error {
	// 定时任务, 每5分钟执行一次, 检查证书
	_, err := cron.Default().AddFunc("*/5 * * * *", crontab.CheckCertHosting)
	if err != nil {
		return err
	}

	handlers := []Handler{
		&handler.CertHandler{},
		&handler.CertHostingHandler{},
	}

	for _, h := range handlers {
		if err := h.Register(ctx); err != nil {
			return err
		}
	}

	return InitTask()
}
