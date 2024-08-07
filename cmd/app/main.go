package main

import (
	"context"
	"flag"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	_ "github.com/codfrm/cago/database/db/sqlite"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/codfrm/cago/pkg/iam"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/iam/audit/audit_db"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/server/cron"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/api"
	"github.com/codfrm/dns-kit/internal/repository/acme_repo"
	"github.com/codfrm/dns-kit/internal/repository/cdn_repo"
	"github.com/codfrm/dns-kit/internal/repository/cert_hosting_repo"
	"github.com/codfrm/dns-kit/internal/repository/cert_repo"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
	"github.com/codfrm/dns-kit/internal/repository/user_repo"
	"github.com/codfrm/dns-kit/internal/task"
	"github.com/codfrm/dns-kit/migrations"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "config file")
	flag.Parse()

	ctx := context.Background()
	cfg, err := configs.NewConfig("dns-kit", configs.WithConfigFile(configFile))
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}

	// 注册储存实例
	user_repo.RegisterUser(user_repo.NewUser())
	domain_repo.RegisterDomain(domain_repo.NewDomain())
	provider_repo.RegisterProvider(provider_repo.NewProvider())
	cert_repo.RegisterCert(cert_repo.NewCert())
	acme_repo.RegisterAcme(acme_repo.NewAcme())
	cdn_repo.RegisterCdn(cdn_repo.NewCdn())
	cert_hosting_repo.RegisterCertHosting(cert_hosting_repo.NewCertHosting())

	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(db.Database()).
		Registry(cago.FuncComponent(broker.Broker)).
		Registry(cache.Cache()).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(db.Default())
		})).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			storage, err := audit_db.NewDatabaseStorage(db.Default())
			if err != nil {
				return err
			}
			return iam.IAM(user_repo.User(), iam.WithAuditOptions(audit.WithStorage(storage)))(ctx, cfg)
		})).
		Registry(cron.Cron()).
		Registry(cago.FuncComponent(task.Task)).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
