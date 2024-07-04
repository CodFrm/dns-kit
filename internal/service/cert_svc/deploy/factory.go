package deploy

import (
	"context"
	"errors"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
)

type Deploy interface {
	Domains(ctx context.Context, hosting *cert_hosting_entity.CertHosting) ([]string, error)
	Deploy(ctx context.Context, hosting *cert_hosting_entity.CertHosting, cert *cert_entity.Cert) error
}

func Factory(ctx context.Context, hosting *cert_hosting_entity.CertHosting) (Deploy, error) {
	provider, err := provider_repo.Provider().Find(ctx, hosting.ProviderID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	switch provider.Platform {
	case provider_entity.PlatformKubernetes:
		return NewK8s(provider), nil
	}
	return nil, errors.New("not support")
}
