package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/db"
	_ "github.com/codfrm/cago/database/db/sqlite"
	"github.com/codfrm/cago/pkg/component"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/api"
	"github.com/codfrm/dns-kit/migrations"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("dns-kit")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	err = cago.New(ctx, cfg).
		Registry(component.Core()).
		Registry(component.Database()).
		Registry(component.Broker()).
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
