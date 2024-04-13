package main

import (
	"context"
	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/db"
	_ "github.com/codfrm/cago/database/db/sqlite"
	"github.com/codfrm/cago/pkg/component"
	"github.com/codfrm/cago/pkg/iam"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/iam/audit/audit_db"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/api"
	"github.com/codfrm/dns-kit/internal/repository/dns_provider_repo"
	"github.com/codfrm/dns-kit/internal/repository/dns_repo"
	"github.com/codfrm/dns-kit/internal/repository/user_repo"
	"github.com/codfrm/dns-kit/migrations"
	"log"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("dns-kit")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}

	// 注册储存实例
	user_repo.RegisterUser(user_repo.NewUser())
	dns_repo.RegisterDns(dns_repo.NewDns())
	dns_provider_repo.RegisterDnsProvider(dns_provider_repo.NewDnsProvider())

	err = cago.New(ctx, cfg).
		Registry(component.Core()).
		Registry(component.Database()).
		Registry(component.Broker()).
		Registry(component.Cache()).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			storage, err := audit_db.NewDatabaseStorage(db.Default())
			if err != nil {
				return err
			}
			return iam.IAM(user_repo.User(), iam.WithAuditOptions(audit.WithStorage(storage)))(ctx, cfg)
		})).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(db.Default())
		})).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
