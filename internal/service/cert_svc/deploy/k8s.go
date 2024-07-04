package deploy

import (
	"context"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"github.com/codfrm/dns-kit/pkg/platform"
	"strings"
)

type K8s struct {
	provider *provider_entity.Provider
}

func NewK8s(provider *provider_entity.Provider) *K8s {
	return &K8s{
		provider: provider,
	}
}

func (k *K8s) Domains(ctx context.Context, hosting *cert_hosting_entity.CertHosting) ([]string, error) {
	return strings.Split(hosting.ConfigMap()["domain"], ","), nil
}

func (k *K8s) Deploy(ctx context.Context, hosting *cert_hosting_entity.CertHosting, cert *cert_entity.Cert) error {
	// 部署到k8s
	manager, err := k.provider.CDNManger(ctx)
	if err != nil {
		return err
	}
	domains := strings.Split(hosting.ConfigMap()["domain"], ",")
	if err := manager.SetCDNHttpsCert(ctx, &platform.CDNItem{
		ID:     hosting.ConfigMap()["namespace"],
		Domain: domains[0],
	}, cert.Certificate, cert.PrivateKey); err != nil {
		return err
	}
	return nil
}
