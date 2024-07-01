package kubernetes

import (
	"context"
	"github.com/codfrm/dns-kit/pkg/platform"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	cli *kubernetes.Clientset
	cfg *rest.Config
}

func NewKubernetes(kubeConfig string) (platform.CDNManager, error) {
	cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return nil, err
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Kubernetes{cli: cli, cfg: cfg}, nil
}

func (k *Kubernetes) GetCDNList(ctx context.Context) ([]*platform.CDNItem, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Kubernetes) GetCDNDetail(ctx context.Context, domain *platform.CDNItem) (*platform.CDNItem, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Kubernetes) SetCDNHttpsCert(ctx context.Context, domain *platform.CDNItem, cert, key string) error {
	//TODO implement me
	panic("implement me")
}

func (k *Kubernetes) UserDetails(ctx context.Context) (*platform.User, error) {
	return &platform.User{
		ID:       k.cfg.Host,
		Username: k.cfg.Host,
	}, nil
}
